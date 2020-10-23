import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Switch } from 'react-router-dom';

import { Provider } from 'react-redux';

import store from './store';
import { Routes } from './constants';
import Products from './routes/products';
import Cart from './routes/cart';
import ShopingRoute from './routes/shopingRoute';

ReactDOM.render(
  (
    <Provider store={store}>
      <Router>
        <Switch>
          <ShopingRoute exact path={Routes.home.path} component={Products} />
          <ShopingRoute path={Routes.cart.path} component={Cart} />
        </Switch>
      </Router>
    </Provider>
  ), document.getElementById('root'),
);
