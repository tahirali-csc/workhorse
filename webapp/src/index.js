import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import { BrowserRouter as Router, Switch, Route, Link } from 'react-router-dom';

import * as serviceWorker from './serviceWorker';
import Header from './Header/Header';
import ProjectList from './ProjectList/ProjectList';
import ProjectStatus from './ProjectStatus/ProjectStatus';
import BuildLogs from './BuildLogs/BuildLogs';

ReactDOM.render(
  <React.StrictMode>
    {/* <App /> */}
    <div>
      <Header />
      <div>
        <Router>
          <Switch>
            <Route exact path='/' component={ProjectList} />
            <Route exact path='/projectStatus/:projectId' component={ProjectStatus} />
            <Route exact path='/buildLogs/:buildId' component={BuildLogs} />
          </Switch>
        </Router>
      </div>
    </div>
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
