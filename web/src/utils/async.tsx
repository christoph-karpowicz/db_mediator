function sleep(ms: number): Promise<Function> {
    return new Promise(resolve => setTimeout(resolve, ms));
}

export { sleep };