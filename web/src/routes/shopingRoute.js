import React from "react";
import { Route } from "react-router-dom";
import { connect } from "react-redux";

const ShopingRoute = ({ component: Component, itemsInCart, ...rest }) => (
  <Route
    {...rest}
    render={(props) => <Component {...props} itemsInCart={itemsInCart} />}
  />
);

function mapStateToProps(state) {
  return {
    itemsInCart: state.products.cartItems.length,
  };
}

export default connect(mapStateToProps, null)(ShopingRoute);
