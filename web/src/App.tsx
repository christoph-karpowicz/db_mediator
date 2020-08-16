import React from 'react';
import { BrowserRouter as Router } from "react-router-dom";
import Application from './app/app';
import './css/App.css';
import Navigation from './componenets/nav/Navigation';
import SubNavigation from './componenets/nav/SubNavigation'
import Content from './componenets/content/Content';

function App(): JSX.Element {
  Application.init();
  
  const [subNavigationActive, setSubNavigationActive] = 
    React.useState(false)
  const [watchers, setWatchers] = React.useState(null);

  function toggleSubNavigationActive(): void {
    setSubNavigationActive(!subNavigationActive);
  }
  
  return (
    <Router>
      <div className="App">
        <Navigation 
          isSubNavigationActive={subNavigationActive} 
          toggleSubNavigationActive={toggleSubNavigationActive} 
          setWatchers={setWatchers}
        />
        <SubNavigation 
          toggleSubNavigationActive={toggleSubNavigationActive} 
          isActive={subNavigationActive}
          watchers={watchers}
        />
        <Content />
      </div>
    </Router>
  );
}

export default App;
