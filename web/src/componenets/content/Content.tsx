import React from 'react';
import '../../css/Content.css';
import { Switch, Route } from "react-router-dom";
import DashboardSection from './sections/dashboard/DashboardSection';
import SynchsSection from './sections/synchs/SynchsSection';

function Content() {
  return (
    <div className="content">
        <Switch>
          <Route path="/synchs/:name">
              <SynchsSection />
          </Route>
          <Route path="/">
              <DashboardSection />
          </Route>
        </Switch>
    </div>
  );
}

export default Content;
