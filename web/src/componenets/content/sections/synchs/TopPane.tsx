import React from 'react';
import { useParams } from "react-router-dom";
import WS from '../../../../ws/ws';
import WSRequest from '../../../../ws/request';

function TopPane() {
    const { name } = useParams();
    
    const onStartSynch = () => {
        const req = new WSRequest("startSynch", { payload: name });
        WS.getSocket().then((ws) => {
            ws.emitAndAwaitResponse(req)
            .then((res: any) => {
                try {
                    const result = JSON.parse(res.Data.Payload);
                } catch(e) {
                    console.error(e);
                }
                console.log(res);
            })
            .catch((err: any) => {
                console.error(err);
            });
        });
    };

    const onStopSynch = () => {
        const req = new WSRequest("stopSynch", { payload: name });
        WS.getSocket().then((ws) => {
            ws.emitAndAwaitResponse(req)
            .then((res: any) => {
                try {
                    const result = JSON.parse(res.Data.Payload);
                } catch(e) {
                    console.error(e);
                }
                console.log(res);
            })
            .catch((err: any) => {
                console.error(err);
            });
        });
    };
    
    return (
        <div className="synchs-top-pane">
            <button className="synchs-top-pane-button" onClick={onStartSynch}>
                Start
            </button>
            <button className="synchs-top-pane-button" onClick={onStopSynch}>
                Stop
            </button>
        </div>
    );
}

export default TopPane;
