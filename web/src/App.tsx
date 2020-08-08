import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import WS from './ws/ws';
import WSReqiest from './ws/request';
import './App.css';
import { ReactComponent as WatchersIcon } from './assets/watchers.svg';
import WatchersSection from './WatchersSection';

function App() {
  const ws = new WS("ws://127.0.0.1:8000/ws/");
  ws.init();

  let req = new WSReqiest("test", {a:1});
  ws.emit(req.json);
  setTimeout(() => {
    ws.emit(req.json);
  }, 4000);
  // while (!ws.emit(req.json));

  return (
    <Router>
      <div className="App">
        <nav>
          <ul>
            <li>
              <Link to="/"><WatchersIcon /></Link>
            </li>
          </ul>
        </nav>
        <div className="content">
          <Switch>
            <Route path="/">
              <WatchersSection />
            </Route>
          </Switch>
        </div>
      </div>
    </Router>
  );
}

export default App;
