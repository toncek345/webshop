import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';

import { Provider } from 'react-redux';

import store from './store';
import { Routes } from './constants';
import Home from './routes/home';

ReactDOM.render(
  (
    <Provider store={store}>
      <Router>
        <Switch>
          <Route path={Routes.home.path} component={Home} />
        </Switch>
      </Router>
    </Provider>
  ), document.getElementById('root'),
);
