class WSRequest {
    private _name: string;
    private _data: object;

    constructor(name: string, data: object) {
        this._name = name;
        this._data = data;
    }

    get json(): string {
        return JSON.stringify({
            name: this._name,
            data: this._data,
        });
    }
}

export default WSRequest;
