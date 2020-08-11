import Application from '../app/app';
import WSRequest from './request';
import { sleep } from '../utils/async';

class WS {
    private static _instance: WS;
    private _connString: string;
    private _socket: WebSocket;
    private _isConnected: boolean;
    private static readonly TIMEOUT_AFTER: number = 2500;

    private constructor(connString: string) {
        this._connString = connString;
    }

    public static async getSocket(): Promise<WS> {
        if (!WS._instance) {
            WS._instance = new WS("ws://127.0.0.1:8000/ws/");
            await WS._instance.init();
        }
        return WS._instance;
    }
    
    public async init(): Promise<void> {
        console.log('Websocket init.')
        this._socket = new WebSocket(this._connString);
        await this.setOnOpen();
        this.setOnClose();
        this.setOnError();
        this.setOnMessage();
    }

    private setOnOpen(): Promise<boolean> {
        return new Promise<boolean>((resolve, reject) => {
            this._socket.onopen = () => {
                console.log("Successfully Connected");
                this.isConnected = true;
                resolve(true)
            };
        });
    }

    private setOnClose() {
        this._socket.onclose = event => {
            console.log("Socket Closed Connection: ", event);
            this.isConnected = false;
            // this.emit("Client Closed!");
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
            try {
                const response = JSON.parse(event.data);
                Application.wsRequestPool.respond(response);
            } catch(e) {
                console.error(e);
            }
        }
    }
    
    public emitAndAwaitResponse(req: WSRequest): Promise<object> {
        req.setExpectResponse(true);
        
        return new Promise<object>((resolve, reject) => {
            const timeStart = new Date().getTime();
            
            const emitted: boolean = this.emit(req);
            if (!emitted) {
                reject({ msg: "Websocket is closed." });
                return;
            }

            const awaitResponse = (initial?: boolean) => {
                const sleepFor = initial ? 1 : 1000;
                
                sleep(sleepFor).then(() => {
                    const currentTime = new Date().getTime();
                    const timeDiff = currentTime - timeStart;
                    // console.log(timeDiff)
                    if (timeDiff > WS.TIMEOUT_AFTER) {
                        reject({ msg: `Request with ID: ${req.getId()} timed out.` });
                        return;
                    }
    
                    if (Application.wsRequestPool.hasResponse(req.getId())) {
                        resolve(Application.wsRequestPool.poll(req.getId()));
                        return;
                    }

                    awaitResponse();
                });
            }

            awaitResponse(true);
        });
    }

    public emit(req: WSRequest): boolean {
        if (!this.isConnected) {
            return false;
        }
        
        this._socket.send(req.json);
        Application.wsRequestPool.append(req);
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
