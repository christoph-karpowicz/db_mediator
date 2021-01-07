import React, { useState } from 'react';
import '../../css/Navigation.css';
import { Link } from "react-router-dom";
import WS from '../../ws/ws';
import WSRequest from '../../ws/request';
import { ReactComponent as DashboardIcon } from '../../assets/dashboard.svg';
import { ReactComponent as WatchersIcon } from '../../assets/watchers.svg';

function Navigation(props: any) {
  function onSynchsClick(): void {
    const req = new WSRequest("getSynchsList", {});
    
    WS.getSocket().then((ws) => {
      ws.emitAndAwaitResponse(req)
        .then((res: any) => {
          try {
            const synchs = JSON.parse(res.Data.Payload);
            props.setSynchs(synchs);
            props.toggleSubNavigationActive();
          } catch(e) {
            console.error(e);
          }
          console.log(res);
        })
        .catch((err: any) => {
          console.error(err);
        });
    });
  }
  
  return (
    <nav className={props.isSubNavigationActive ? "" : "bordered"}>
        <ul>
            <li>
                <Link to="/"><DashboardIcon /></Link>
            </li>
            <li>
                <WatchersIcon onClick={onSynchsClick} />
            </li>
        </ul>
    </nav>
  );
}

export default Navigation;
