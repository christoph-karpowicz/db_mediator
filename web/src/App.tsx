import React from 'react';
import { BrowserRouter as Router } from "react-router-dom";
import WS from './ws/ws';
import WSRequest from './ws/request';
import './css/App.css';
import Navigation from './componenets/nav/Navigation';
import SubNavigation from './componenets/nav/SubNavigation'
import Content from './componenets/content/Content';

function App(): JSX.Element {
  const [subNavigationActive, setSubNavigationActive] = 
    React.useState(false)
  
  // const ws = new WS("ws://127.0.0.1:8000/ws/");
  // ws.init();

  // let req = new WSRequest("test", {a:1});
  // ws.emit(req.json);
  // setTimeout(() => {
  //   ws.emit(req.json);
  // }, 4000);
  // while (!ws.emit(req.json));

  function onNavClick(): void {
    console.log('ss')
    setSubNavigationActive(!subNavigationActive);
  }
  
  return (
    <Router>
      <div className="App">
        <Navigation 
          isSubNavigationActive={subNavigationActive} 
          onNavClick={onNavClick} 
        />
        <SubNavigation 
          isActive={subNavigationActive} 
        />
        <Content />
      </div>
    </Router>
  );
}

export default App;
