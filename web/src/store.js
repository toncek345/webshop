import { createStore, applyMiddleware, combineReducers } from "redux";
import logger from "redux-logger";
import thunk from "redux-thunk";

import { composeWithDevTools } from "redux-devtools-extension";

import Products from "./reducers/products";

const middleware = composeWithDevTools(applyMiddleware(thunk, logger));

const reducers = combineReducers({
  products: Products,
});

export default createStore(reducers, middleware);
