package main

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/packr"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type DirFiles struct {
	Path string `json:",omitempty"`
	Dir  bool `json:",omitempty"`
	Uuid string `json:",omitempty"`
	List []DirFiles
	Name string
}
func hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func getFiles(root string) []DirFiles {
	files := make([]DirFiles, 0)
	file1, _ := ioutil.ReadDir(root)
	for _, f := range file1 {
		if f.IsDir(){
			files = append(files, DirFiles{Path: root+"/"+f.Name(), Dir: f.IsDir(),
				Uuid: hash(root+"/"+f.Name()), List:getFiles(root+"/"+f.Name()), Name:f.Name()})
		}else{
			files = append(files, DirFiles{Path: root+"/"+f.Name(), Dir: f.IsDir(),
				Uuid: hash(root+"/"+f.Name()), List:make([]DirFiles, 0), Name:f.Name()})
		}

	}
	return files
}


func fileHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["path"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}
	http.ServeFile(w, r, keys[0])
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["path"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}
	Openfile, err := os.Open(keys[0])
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		http.Error(w, "File not found.", 404)
		return
	}

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string
	dir, file := filepath.Split(keys[0])
	fmt.Println("Dir:", dir)   //Dir: /some/path/to/remove/
	fmt.Println("File:", file) //File: ile.name
	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename=\""+file+"\"")
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(w, Openfile) //'Copy' the file to the client
	return
}

func dirHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["path"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}
	res1B, _ := json.Marshal(getFiles(keys[0]))
	fmt.Fprintln(w, string(res1B))
}
func main() {
	key := [256]byte{}
	_, err := rand.Read(key[:])
	if err != nil {
		panic(err)
	}
	box := packr.NewBox("./dist")
	http.Handle("/", http.FileServer(box))
	http.HandleFunc("/dir", dirHandler)
	http.HandleFunc("/f", fileHandler)
	http.HandleFunc("/d", downloadHandler)
	panic(http.ListenAndServe(":8080", nil))
}
