import React from 'react';
import '../../css/Content.css';
import { Switch, Route } from "react-router-dom";
import WatchersSection from './WatchersSection';

function Content() {
  return (
    <div className="content">
        <Switch>
        <Route path="/">
            <WatchersSection />
        </Route>
        </Switch>
    </div>
  );
}

export default Content;
