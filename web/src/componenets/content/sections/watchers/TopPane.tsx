import React from 'react';
import { useParams } from "react-router-dom";
import WS from '../../../../ws/ws';
import WSRequest from '../../../../ws/request';

function TopPane() {
    const { name } = useParams();
    
    const onStartWatcher = () => {
        const req = new WSRequest("startWatcher", { payload: name });
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

    const onStopWatcher = () => {
        const req = new WSRequest("stopWatcher", { payload: name });
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
        <div className="watchers-top-pane">
            <button className="watchers-top-pane-button" onClick={onStartWatcher}>
                Start
            </button>
            <button className="watchers-top-pane-button" onClick={onStopWatcher}>
                Stop
            </button>
        </div>
    );
}

export default TopPane;
