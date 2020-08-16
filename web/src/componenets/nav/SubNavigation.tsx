import React, { useRef, useState, useEffect } from 'react';
import '../../css/SubNavigation.css';
import { Link } from "react-router-dom";

function SubNavigation(props: any) {
  const subnav = useRef<HTMLDivElement>(null);
  const [left, setLeft] = useState<number>(0);

  useEffect(() => {
    if (subnav.current) {
      const current: { offsetWidth: number } = subnav.current;
      const width: number = current.offsetWidth;
      const left = width > 0 ? width : 200;
      setLeft(left);
    }
  })
  
  return (
    <div 
      id="sub-navigation" 
      onClick={props.toggleSubNavigationActive}
      style={{ left: (props.isActive ? left : -1*left) + "px" }}
      ref={subnav}
    >
        <ul>
            {props.watchers && props.watchers.map((watcher: string, i: number) => {
              return (
                <li key={i}>
                    <Link to={"/watchers/" + watcher}>{watcher}</Link>
                </li>
              );
            })}
        </ul>
    </div>
  );
}

export default SubNavigation;
