export function delay(milliseconds: number) {
    return new Promise(function run(resolve) {
        setTimeout(resolve, milliseconds);
    });
}

