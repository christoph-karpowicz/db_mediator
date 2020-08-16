import React from 'react';
import '../../css/Content.css';
import { Switch, Route } from "react-router-dom";
import DashboardSection from './sections/dashboard/DashboardSection';
import WatchersSection from './sections/watchers/WatchersSection';

function Content() {
  return (
    <div className="content">
        <Switch>
          <Route path="/watchers/:name">
              <WatchersSection />
          </Route>
          <Route path="/">
              <DashboardSection />
          </Route>
        </Switch>
    </div>
  );
}

export default Content;
