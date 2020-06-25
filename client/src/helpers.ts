export const copyStringToClipboard = (str: any): void => {
    const el = document.createElement('textarea') as HTMLTextAreaElement;
    el.value = str;
    el.setAttribute('readonly', '');
    el.style.position = 'absolute';
    el.style.left = '-9999px';
    document.body.appendChild(el);
    el.select();
    document.execCommand('copy');
    document.body.removeChild(el);
};

export const shuffle = (array: any): any[] => {
    let currentIndex = array.length;
    let temporaryValue = null;
    let randomIndex = null;
    while (0 !== currentIndex) {
        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex -= 1;
        temporaryValue = array[currentIndex];
        array[currentIndex] = array[randomIndex];
        array[randomIndex] = temporaryValue;
    }

    return array;
};

export const getExtension = (file: string): string => {
    if(file === undefined){
        return '';
    }
    const basename = file.split(/[\\/]/).pop();
    if (basename !== undefined) {
        const pos = basename.lastIndexOf('.');

        if (basename === '' || pos < 1) {
            return '';
        }
        return basename.slice(pos + 1);
    }
    return '';

};