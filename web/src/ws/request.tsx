import { uuid } from 'uuidv4';

class WSRequest {
    private _id: string;
    private _name: string;
    private _data: object;
    private _expectResponse: boolean;

    constructor(name: string, data: object) {
        this._name = name;
        this._data = data;
    }

    private assignId(): void {
        this._id = uuid();
    }

    get json(): string {
        let req: object;
        if (this._expectResponse) {
            this.assignId();
        }

        if (this._data) {
            req = {
                id: this._id,
                name: this._name,
                data: this._data
            };
        } else {
            req = { 
                id: this._id, 
                name: this._name 
            };
        }
        
        return JSON.stringify(req);
    }

    public getId(): string {
        return this._id;
    }

    public setExpectResponse(expectResponse: boolean): void {
        this._expectResponse = expectResponse;
    }
}

export default WSRequest;
