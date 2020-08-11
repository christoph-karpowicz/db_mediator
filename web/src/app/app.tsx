import WSRequestPool from '../ws/request_pool';

class Application {
    public static wsRequestPool: WSRequestPool;

    public static init() {
        this.wsRequestPool = new WSRequestPool();
    }

}

export default Application;
