class WS {
    private _connString: string;
    private _socket: WebSocket;
    private _isConnected: boolean;

    constructor(connString: string) {
        this._connString = connString;
    }
    
    public init() {
        this._socket = new WebSocket(this._connString);
        this.setOnOpen();
        this.setOnClose();
        this.setOnError();
        this.setOnMessage();
    }

    private setOnOpen() {
        this._socket.onopen = () => {
            console.log("Successfully Connected");
            this.isConnected = true;
            this._socket.send("test")
        };
    }

    private setOnClose() {
        this._socket.onclose = event => {
            console.log("Socket Closed Connection: ", event);
            this.isConnected = false;
            this.emit("Client Closed!");
        };
    }

    private setOnError() {
        this._socket.onerror = error => {
            console.log("Socket Error: ", error);
            this.isConnected = false;
        };
    }

    private setOnMessage() {
        this._socket.onmessage = function (event) {
            console.log(event.data);
        }
    }
    
    public emit(msg: string): boolean {
        if (!this.isConnected) {
            console.error("Websocket is closed.");
            return false;
        }
        
        this._socket.send(msg);
        return true;
    }

    set isConnected(isConnected: boolean) {
        this._isConnected = isConnected;
    }

    get isConnected(): boolean {
        return this._isConnected;
    }
    
    get socket(): WebSocket {
        return this._socket;
    }
}

export default WS;
