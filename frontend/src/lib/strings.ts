export function capitalize(str: string, delimiters: string[] = [' ', '/']): string {
    if (!str) return '';
    for (const delimiter of delimiters) {
        str = str.replaceAll(delimiter, ' ');
    }
    const words: string[] = str
        .split(' ')
        .map((word) => `${word[0].toUpperCase()}${word.slice(1)}`);
    /*
        for(const [index, word]	of words){
        words[index] = word[0].toUpperCase()
    }
    */
    return words.join(' ');
}