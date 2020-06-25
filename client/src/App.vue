<template>
    <v-app light>

        <v-navigation-drawer :clipped="clipped" v-model="drawer" app
                             :right="true" :disable-resize-watcher="true" width="400">
            <v-list>
                <v-list-tile
                        v-for="(track, index) in playlist"
                        :key="index"
                        avatar
                        @click="select(index)"
                >
                    <v-list-tile-action>
                        <v-icon v-if="index == currentTrackIndex" color="black">music_video</v-icon>
                    </v-list-tile-action>
                    <v-list-tile-content>
                        <v-tooltip bottom>
                            <span slot="activator"> <v-list-tile-title v-text="track.name"></v-list-tile-title></span>
                            <span>{{track.name}}</span>
                        </v-tooltip>
                    </v-list-tile-content>
                </v-list-tile>
            </v-list>
        </v-navigation-drawer>
        <v-toolbar fixed app :clipped-left="clipped">
            <v-toolbar-title>Good Player(by <a href="https://gitlab.com/indigogodman">Максимка</a>)</v-toolbar-title>
            <v-spacer></v-spacer>
            Playlist
            <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
        </v-toolbar>
        <v-content>
            <v-container grid-list-xl xs6>
                <v-layout align-center justify-center row fill-height>
                    <v-flex class="text-xs-center" xs12 sm10 md8 lg8 x8>
                        <v-expansion-panel dark>
                            <v-expansion-panel-content>
                                <div slot="header">Files</div>
                                <v-treeview
                                        ref="tree"
                                        :open.sync="open"
                                        v-model="tree"
                                        :load-children="getTree"
                                        :items="items"
                                        activatable
                                        active-class="grey"
                                        selected-color="white"
                                        open-on-click
                                        selectable
                                        dark
                                >
                                    <template slot="prepend" slot-scope="{ item, open, leaf }" >

                                        <v-icon @click="runCurrent(item)" v-if="(audioTypes.includes(getExtension(item.path)) || videoTypes.includes(getExtension(item.path))) && !compare_src(item.path)">
                                            play_arrow
                                        </v-icon>

                                        <v-icon @click="pause_click(item)" v-if="compare_src(item.path) && !played">
                                            play_arrow
                                        </v-icon>

                                        <v-icon @click="pause_click(item)" v-if="compare_src(item.path) && played">
                                            pause
                                        </v-icon>
                                        <v-icon @click="runCurrentDir(item)" v-if="!audioTypes.includes(getExtension(item.path)) && !videoTypes.includes(getExtension(item.path)) ">
                                            play_arrow
                                        </v-icon>
                                    </template>
                                </v-treeview>

                            </v-expansion-panel-content>
                        </v-expansion-panel>
                    </v-flex>
                </v-layout>
                <v-layout align-center justify-center row fill-height dark>

                    <v-flex class="text-xs-center" xs12 sm10 md8 lg8 x8>
                        <audio v-if="audioTypes.includes(getExtension(src))" :src="src" ref='audio' id="audio"
                               style="display:none;width: 100%!important;height:auto!important;"/>
                        <video v-if="videoTypes.includes(getExtension(src))" :src="src" ref='video' id="video" controls
                               style="width: 100%!important;height:auto!important;">

                        </video>
                        <v-flex>
                            <v-tooltip right>
                                <span slot="activator">
                                 <v-slider v-if="audioTypes.includes(getExtension(src)) || videoTypes.includes(getExtension(src))"
                                           v-model="slider"
                                           step="0.000000000000000000000001"
                                           @click="sliderChange"
                                           append-icon="fast_forward"
                                           prepend-icon="fast_rewind"
                                           color="black"
                                           thumb-label="always"
                                           v-on:start="sliderRemoveListeners"
                                           v-on:end="sliderStartListeners"
                                 >
                                <template
                                        slot="thumb-label"
                                        slot-scope="props"
                                >
          <span>
            {{ Math.floor(slider) }}
          </span>
                                </template>
                            </v-slider>
                                </span>
                                <span>Rewind(a)/Forward(d)</span>
                            </v-tooltip>

                            <v-tooltip right>
                                <span slot="activator">
                                <v-slider v-if="audioTypes.includes(getExtension(src)) || videoTypes.includes(getExtension(src))"
                                          v-model="volume"
                                          append-icon="volume_up"
                                          prepend-icon="volume_down"
                                          color="black"
                                          thumb-label
                                ></v-slider>
                                </span>
                                <span>Down(c)/Up(v)</span>
                            </v-tooltip>

                        </v-flex>
                        <div class="text-xs-center">
                            <v-tooltip bottom>
                                <span slot="activator">
                                    <v-btn outline large fab color="black" @click="previous_track">
                                        <v-icon>skip_previous</v-icon>
                                    </v-btn></span>
                                <span>Previous track (q)</span>
                            </v-tooltip>
                            <v-tooltip bottom>
                                <span slot="activator"> <v-btn outline large fab color="black" @click="play_pause">
                                <v-icon v-if="!played">play_arrow</v-icon>
                                <v-icon v-if="played">pause</v-icon>
                            </v-btn></span>
                                <span>Play/Pause (s)</span>
                            </v-tooltip>
                            <v-tooltip bottom>
                                <span slot="activator"> <v-btn outline large fab color="black" @click="next_track">
                                <v-icon>skip_next</v-icon>
                            </v-btn></span>
                                <span>Next track (e)</span>
                            </v-tooltip>

                            <v-tooltip bottom>
                                <span slot="activator"> <v-btn outline large fab color="black"
                                                               @click="checkbox=!checkbox">
                                <v-icon>shuffle</v-icon>
                            </v-btn></span>
                                <span>Shuffle</span>
                            </v-tooltip>

                            <v-tooltip bottom>
                                <span slot="activator"> <v-btn outline large fab :color="(replay) ? 'indigo': 'black' "
                                                               @click="replay=!replay">
                                <v-icon>replay</v-icon>
                            </v-btn></span>
                                <span>Replay</span>
                            </v-tooltip>

                            <v-tooltip bottom v-if="isVideo">
                                <span slot="activator"> <v-btn outline large fab color="black" @click="fulltoggle">
                                <v-icon>fullscreen</v-icon>
                            </v-btn></span>
                                <span>FullSize (space)</span>
                            </v-tooltip>

                            <v-tooltip bottom>
                                <span slot="activator"> <v-btn outline large fab :color="'black'"
                                                               @click="copytoClipboard">
                                <v-icon>link</v-icon>
                            </v-btn></span>
                                <span>Link (l)</span>
                            </v-tooltip>

                            <v-tooltip bottom>
                                <span slot="activator"> <v-btn outline large fab :color="'black'"
                                                               @click="downloadCurrentFile">
                                <v-icon>cloud_download</v-icon>
                            </v-btn></span>
                                <span>Download</span>
                            </v-tooltip>
                        </div>
                    </v-flex>
                </v-layout>
            </v-container>
        </v-content>
    </v-app>

</template>

<script lang="ts">
    import {Component, Vue, Watch} from 'vue-property-decorator';
    import {copyStringToClipboard, getExtension, shuffle} from "./helpers";

    interface MyVideoElement extends HTMLVideoElement {
        mozFullScreenElement?: boolean;
        webkitFullscreenElement?: boolean;
        webkitRequestFullscreen?: () => void;
        mozRequestFullScreen?: () => void;
        msRequestFullscreen?: () => void;
    }


    @Component
    export default class App extends Vue {
        private api: string = `f?path=`;
        private breweries: any = [];
        private tree: any = [];
        private types: any = [];
        private currentTrackIndex: any = 0;
        private checkbox: boolean = false;
        private src: any = '';
        private predhear:boolean =  false;
        private played: any = false;
        private open: any = [];
        private drawer: any = false;
        private clipped: any = false;
        private volume: any = 50;
        private slider: any = 0;
        private previ: any = false;
        private pl: any = false;
        private replay: any = false;
        private items: any = [{id: '', name: 'media', children: [], isDir: true}];
        private playlist: any[] = [];
        private audioTypes: string[] = ['mp3', 'flac'];
        private videoTypes: string[] = ['mp4'];
        private document: any = window.document;


        private getExtension(need: string){
            return getExtension(need);
        }

        private compare_src(item:any){
            return this.src === this.geturl`${item}`
        }
        /**
         * Extension select file in playlist
         * @returns {string}
         */
        get currentFileExtension(): string {
            return getExtension(this.src);
        }

        /**
         * Is current player
         * @returns {HTMLAudioElement | MyVideoElement}
         */
        get currentPlayer(): HTMLAudioElement | MyVideoElement {
            return this.isAudio ? this.$refs.audio as HTMLAudioElement : this.$refs.video as MyVideoElement;
        }

        private getCurrentPlayer(path: string): HTMLAudioElement | MyVideoElement {
            return this.audioTypes.includes(getExtension(path)) ? this.$refs.audio as HTMLAudioElement : this.$refs.video as MyVideoElement;
        }

        /**
         * Is audio current file
         * @returns {boolean}
         */
        get isAudio(): boolean {
            console.log()
            return this.audioTypes.includes(this.currentFileExtension);
        }

        get duration(){
            if (this.currentPlayer !== undefined) {
                return this.currentPlayer.duration;
            }
            else{
                return 0;
            }

        }

        get currentTime(){
            if (this.currentPlayer !== undefined) {
                return this.currentPlayer.currentTime;
            }
            else{
                return 0;
            }
        }

        set currentTime(value){
            if (this.currentPlayer !== undefined) {
                this.currentPlayer.currentTime = value;
            }
        }


        /**
         * Is video current file
         * @returns {boolean}
         */
        get isVideo(): boolean {
            return this.videoTypes.includes(this.currentFileExtension);
        }

        private runCurrent(item:any){
            this.src = this.geturl`${item.path}`;
            Vue.nextTick(() => {
                if ( this.currentPlayer.onended === null) {
                    this.setListeners()
                }
                this.predhear = true;
                this.currentPlayer.load();
                this.currentPlayer.play().then(() => this.currentPlayer.volume = this.volume / 100);
            });
        }

        private runCurrentDir(item:any){
            const tree = this.$refs.tree as any;
            this.tree = [];
            Object.keys(tree.nodes).forEach(key => tree.updateSelected(key, false))
            tree.updateSelected(item.id, true);
            Vue.nextTick(() => {
                    [...tree.selectedCache].map(x => x ? this.tree.push(x) : '')
                    Vue.nextTick(() => {
                        this.play_pause();
                    });
                }
            )
        }


        private createPlaylist(): string[] {
            if (!this.items.length) {
                return [];
            }
            const playlist = [];

            const tree = this.items[0];

            for (const leaf of this.tree) {
                const found = this._findInTree(leaf, tree);
                if (found !== null) {
                    if (this.supportType(found.result.name)) {
                        playlist.push(found.result);
                    }
                }
            }
            if (this.checkbox) {
                return shuffle(playlist);
            }
            return playlist;
        }

        /**
         *  Emit from created
         *  Create playlist if link is exist in parameters
         *  Event for space key
         */
        private async created(): Promise<void> {
            await this.getTree();
            const uri = window.location.href.split('?');
            if (uri.length === 2) {
                const vars = uri[1].split('&');
                const getVars: any = {};
                let tmp = '';
                vars.forEach((v: any) => {
                    tmp = v.split('=');
                    if (tmp.length === 2) {
                        getVars[tmp[0]] = tmp[1];
                    }
                });
                if (getVars.link !== '') {
                    Vue.nextTick(() => {
                        const tree = this.$refs.tree as any;
                        tree.updateSelected(getVars.link, true);
                        Vue.nextTick(() => {
                            this.tree.push(getVars.link);
                            Vue.nextTick(() => {
                                this.play_pause();
                            });
                        });
                    });
                }
            }
            window.addEventListener('keypress', (event: any) => {
                switch (event.key) {
                    case ' ':
                        this.fulltoggle();
                        break;
                    case 'q':
                        this.previous_track();
                        break;
                    case 'e':
                        this.next_track();
                        break;
                    case 's':
                        this.play_pause();
                        break;
                    case 'a':
                        this.currentTime = this.lowArrow() ? this.lowArrow() : 0;
                        break;
                    case 'd':
                        this.currentTime = this.highArrow() < this.duration ? this.highArrow() : this.duration;
                        break;
                    case 'c':
                        this.volume = this.lowSound()>0 ? this.lowSound() : 0;
                        break;
                    case 'v':
                        this.volume = this.highSound() > 100 ? 100 : this.highSound();
                        break;
                    case 'r':
                        this.replay = !this.replay;
                        break;
                    case 'l':
                        this.copytoClipboard();
                        break;
                }
            })
        }


        /**
         * Download current media file
         * Create link-element and click him
         */
        private downloadCurrentFile(): void {
            const link = document.createElement('a');
            link.setAttribute('href', this.src.replace(/f\?/g, 'd?'));
            link.setAttribute('download', this.src);
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        }

        /**
         * Check support extension for web
         * @param file Checking file
         * @return {boolean}
         */
        private supportType(file: string): boolean {
            const ext = getExtension(file);
            return this.audioTypes.includes(ext) || this.videoTypes.includes(ext);
        }

        private getChildren(array: any) {
            return array.map((x: any) => {
                if (x.Dir) {
                    return {id: x.Uuid, name: x.Name, path: x.Path, children: this.getChildren(x.List), isDir: x.Dir};
                } else {
                    return {id: x.Uuid, name: x.Name, path: x.Path, isDir: x.Dir};
                }
            });
        }

        private lowSound():number{
            if (this.currentPlayer !== undefined) {
                return this.volume - 10;
            }
            return 0;
        }

        private highSound():number{
            if (this.currentPlayer !== undefined) {
                return this.volume + 10;
            }
            return 0;
        }

        private highArrow():number{
            if (this.currentPlayer !== undefined) {
                return this.currentPlayer.currentTime + 10;
            }
            return 0;
        }

        private lowArrow(): number{
            if (this.currentPlayer !== undefined) {
                return this.currentPlayer.currentTime - 10;
            }
            return 0;
        }


        private geturl(str: TemplateStringsArray, val: string) {
            return `${this.api}${val}`;
        }

        private copytoClipboard() {
            if (this.playlist.length) {
                copyStringToClipboard(`${window.location.href}?link=${this.playlist[this.currentTrackIndex].id}`);
            }
        }

        private ended(){
            if (this.playlist.length) {
                if(this.currentTrackIndex == -1){
                    this.playlist[0].active = false;
                }else{
                    this.playlist[this.currentTrackIndex].active = false;
                }

                if (this.previ) {
                    if (this.currentTrackIndex > 0) {
                        this.currentTrackIndex -= 2;
                    }
                    this.previ = false;
                }
                if(this.predhear){
                    this.predhear=false;
                    this.currentTrackIndex-=1;
                }
                if (!this.replay) {
                    this.currentTrackIndex += 1;
                }
                if (this.currentTrackIndex < this.playlist.length) {
                    this.playlist[this.currentTrackIndex].active = true;
                    this.src = this.geturl`${this.playlist[this.currentTrackIndex].path}`;
                    //history.pushState(null, null, `${window.location.href}?link=${this.playlist[this.currentTrackIndex].id}`);
                    Vue.nextTick(() => {
                        if (this.currentPlayer.onended === null) {
                            this.setListeners();
                        }
                        this.currentPlayer.load();
                        this.currentPlayer.play().then(() => this.currentPlayer.volume = this.volume / 100);
                    });
                }
                else {
                    this.currentTrackIndex = 0;
                }
            }

        }

        private select(item: any) {
            this.currentTrackIndex = item - 1;
            Vue.nextTick(() => this.ended());
        }


        private runPlaylist() {
            this.currentTrackIndex = 0;
            this.playlist = this.createPlaylist();
            if (this.playlist.length) {
                this.src = this.geturl`${this.playlist[this.currentTrackIndex].path}`;
                Vue.nextTick(() => {
                    this.playlist[this.currentTrackIndex].active = true;
                    if (this.currentTrackIndex < this.playlist.length) {
                        this.src = this.geturl`${this.playlist[this.currentTrackIndex].path}`;
                        Vue.nextTick(async () => {
                            this.setListeners();
                            this.currentPlayer.load();
                            await this.currentPlayer.play();
                            this.currentPlayer.volume = this.volume / 100;
                        });
                    }
                });
            }
        }

        private play_pause() {
            if (this.pl) {
                this.pl = false;
                this.runPlaylist();
            } else {
                if (this.src!=="") {
                    Vue.nextTick(async () => {
                        if (!this.played) {
                            await this.currentPlayer.play();
                            this.currentPlayer.volume = this.volume / 100;
                        } else {
                            this.currentPlayer.pause();
                        }
                    });
                }
            }
        }

        private pause_click(){
            if (this.src!=="") {
                Vue.nextTick(async () => {
                    if (!this.played) {
                        await this.currentPlayer.play();
                        this.currentPlayer.volume = this.volume / 100;
                    } else {
                        this.currentPlayer.pause();
                    }
                });
            }
        }


        private setListeners() {
            this.currentPlayer.onended = this.ended.bind(this);
            this.currentPlayer.onplaying = () => this.played = true;
            this.currentPlayer.onpause = () => this.played = false;
            this.sliderStartListeners(false);
            this.currentPlayer.onerror = this.next_track.bind(this);
        }


        private sliderStartListeners(slider=true){
            if(slider){
                this.sliderChange();
            }
            this.currentPlayer.ontimeupdate = () => {
                const newTime = this.currentPlayer.currentTime / this.currentPlayer.duration * 100;
                if (newTime) {
                    this.slider = newTime;
                }
            };
        }


        private sliderRemoveListeners(){
            this.currentPlayer.ontimeupdate = () => '';
        }

        private sliderChange() {
            this.currentPlayer.currentTime = this.currentPlayer.duration * (this.slider / 100);
        }

        @Watch('breweries')
        private onBreweriesChange(val: any) {
            this.types = val.reduce((acc: any, cur: any) => {
                const type = cur.brewery_type;

                if (!acc.includes(type)) {
                    acc.push(type);
                }

                return acc;
            }, []).sort();
        }

        private fulltoggle() {

            if (this.playlist.length > 0 && this.isVideo) {
                console.log(this.playlist)
                const element = document.getElementById('video') as MyVideoElement;
                if (!this.document.fullscreenElement && !this.document.mozFullScreenElement &&
                    !this.document.webkitFullscreenElement && !this.document.msFullscreenElement) {
                    if (element.requestFullscreen) {
                        element.requestFullscreen();
                    } else if (element.msRequestFullscreen) {
                        element.msRequestFullscreen();
                    } else if (element.mozRequestFullScreen) {
                        element.mozRequestFullScreen();
                    } else if (element.webkitRequestFullscreen) {
                        element.webkitRequestFullscreen();
                    }
                } else {
                    if (this.document.exitFullscreen) {
                        this.document.exitFullscreen();
                    } else if (this.document.msExitFullscreen) {
                        this.document.msExitFullscreen();
                    } else if (this.document.mozCancelFullScreen) {
                        this.document.mozCancelFullScreen();
                    } else if (this.document.webkitExitFullscreen) {
                        this.document.webkitExitFullscreen();
                    }
                }
            }
        }

        private next_track() {
            this.ended();
        }

        private previous_track() {
            this.previ = true;
            Vue.nextTick(() => this.ended());
        }

        private _findInTree(id: any, tree: any): any {
            let result: any;
            const index: any = null;
            const rv: any = null;
            if (tree.id === id) {
                result = tree;
            } else {
                if (tree.hasOwnProperty('children')) {
                    for (const i of tree.children) {
                        const v = this._findInTree(id, i);
                        if (v != null) {
                            return v;
                        }
                    }
                }
                return null;
            }
            return {result};
        }

        private async getTree() {
            if (this.breweries.length) {
                return;
            }
            const res = await fetch('/dir?path=./media');
            this.breweries = await res.json();
            this.items = [{id: '', name: 'media', children: this.getChildren(this.breweries), isDir: true}];
        }

        @Watch('tree')
        private onchangeTree() {
            this.pl = true;
        }

        @Watch('volume')
        private volumeChange() {
            this.currentPlayer.volume = this.volume / 100;
        }

    }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
    .v-slider__thumb-container v-slider__thumb-container--is-active primary--text {
        background: blue;
    }

    body, html {
        margin: 0;
        padding: 0;
        height: 100%;
    }

    body {
        background-size: cover;
        font-family: 'Cabin Condensed', sans-serif;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
    }

    svg {
        font-weight: bold;
        max-width: 600px;
        height: auto;
    }
</style>
