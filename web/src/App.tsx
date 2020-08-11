import React from 'react';
import { BrowserRouter as Router } from "react-router-dom";
import Application from './app/app';
import WS from './ws/ws';
import WSRequest from './ws/request';
import './css/App.css';
import Navigation from './componenets/nav/Navigation';
import SubNavigation from './componenets/nav/SubNavigation'
import Content from './componenets/content/Content';

function App(): JSX.Element {
  Application.init();
  
  const [subNavigationActive, setSubNavigationActive] = 
    React.useState(false)
  

  // let req = new WSRequest("test", {a:1});
  // ws.emit(req.json);
  // setTimeout(() => {
  //   ws.emit(req.json);
  // }, 4000);
  // while (!ws.emit(req.json));

  function toggleSubNavigationActive(): void {
    setSubNavigationActive(!subNavigationActive);
  }
  
  return (
    <Router>
      <div className="App">
        <Navigation 
          isSubNavigationActive={subNavigationActive} 
          toggleSubNavigationActive={toggleSubNavigationActive} 
        />
        <SubNavigation 
          toggleSubNavigationActive={toggleSubNavigationActive} 
          isActive={subNavigationActive} 
        />
        <Content />
      </div>
    </Router>
  );
}

export default App;
