import React from 'react';
import { Route, HashRouter } from 'react-router-dom';
import Home from './routes/home/Home';
import Host from './routes/host/Host';
import Worker from './routes/worker/Worker';
import './App.css';

const App = () => {

  return (
    <HashRouter>
      <Route path="/" exact={true} component={Home} />
      <Route path="/host" component={Host} />
      <Route path="/worker" component={Worker} />
    </HashRouter>
  );
}

export default App;
