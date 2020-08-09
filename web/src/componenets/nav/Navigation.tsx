import React from 'react';
import '../../css/Navigation.css';
import { Link } from "react-router-dom";
import WS from '../../ws/ws';
import WSRequest from '../../ws/request';
import { ReactComponent as WatchersIcon } from '../../assets/watchers.svg';

function Navigation(props: any) {
  function onWatchersClick(): void {
    let req = new WSRequest("getWatchersList", {});
    WS.getSocket().emit(req.json);

    // console.log('ss')
    props.toggleSubNavigationActive();
  }
  
  return (
    <nav className={props.isSubNavigationActive ? "" : "bordered"}>
        <ul>
            <li>
                {/* <Link to="/"><WatchersIcon /></Link> */}
                <WatchersIcon onClick={onWatchersClick} />
            </li>
        </ul>
    </nav>
  );
}

export default Navigation;
