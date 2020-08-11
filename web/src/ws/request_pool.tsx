import WSRequest from './request'

interface PoolObject {
    [key: string]: any
}

class WSRequestPool {
    private _pool: PoolObject;

    constructor() {
        this._pool = {};
    }

    public append(req: WSRequest): void {
        this._pool[req.getId()] = null;
    }

    public respond(response: { ID: string }): void {
        this._pool[response.ID] = response;
    }

    public hasResponse(reqId: string): boolean {
        return this._pool[reqId] != null;
    }

    public poll(reqId: string): object {
        const response = this._pool[reqId];
        delete this._pool[reqId];
        return response;
    }

}

export default WSRequestPool;
