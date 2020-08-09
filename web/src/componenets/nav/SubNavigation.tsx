import React from 'react';
import '../../css/SubNavigation.css';
import { Link } from "react-router-dom";

function SubNavigation(props: any) {
  return (
    <div id="sub-navigation" className={props.isActive ? "active" : ""} onClick={props.toggleSubNavigationActive}>
        <ul>
            <li>
                <Link to="/">test1</Link>
            </li>
            <li>
                <Link to="/">test2</Link>
            </li>
        </ul>
    </div>
  );
}

export default SubNavigation;
