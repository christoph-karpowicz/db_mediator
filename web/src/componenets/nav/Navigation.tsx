import React from 'react';
import '../../css/Navigation.css';
import { Link } from "react-router-dom";
import { ReactComponent as WatchersIcon } from '../../assets/watchers.svg';

function Navigation(props: any) {
  return (
    <nav className={props.isSubNavigationActive ? "" : "bordered"}>
        <ul>
            <li>
                {/* <Link to="/"><WatchersIcon /></Link> */}
                <WatchersIcon onClick={props.onNavClick} />
            </li>
        </ul>
    </nav>
  );
}

export default Navigation;
