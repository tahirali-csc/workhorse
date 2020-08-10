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
import AccordList from './AccordList'

ReactDOM.render(
  
    <div>
      <Header />
      <div>
        <Router>
          <Switch>
            <Route exact path='/' component={ProjectList} />
            <Route exact path='/projectStatus/:projectId/:name' component={ProjectStatus} />
            <Route exact path='/buildLogs/:buildId' component={BuildLogs} />
            {/* <Route exact path='/buildLogs/:buildId' component={AccordList} /> */}
          </Switch>
        </Router>
      </div>
    </div>
  ,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
