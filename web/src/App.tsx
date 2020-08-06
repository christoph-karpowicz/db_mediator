// import socketIO from 'socket.io-client';
import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import './App.css';
import { ReactComponent as WatchersIcon } from './assets/watchers.svg';
import WatchersSection from './WatchersSection';

function App() {
  // const socket = socketIO('http://localhost:8000', { forceNew: true, path: '/ws' });

  // // function subscribeToTimer(cb) {
  //   socket.on('timer', (m: any) => {
  //     console.log(m)
  //   });
  //   socket.emit('subscribeToTimer', 1000);
  // // }
  
  let socket = new WebSocket("ws://127.0.0.1:8000/ws/");
  console.log("Attempting Connection...");

  socket.onopen = () => {
    console.log("Successfully Connected");
    socket.send("test")
  };
  
  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
    socket.send("Client Closed!")
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };

  socket.onmessage = function (event) {
    console.log(event.data);
  }
  
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
