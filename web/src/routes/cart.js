import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { Grid, Row } from 'react-bootstrap';

import Navigation from '../components/Navigation';
import CartCard from '../components/CartCard';

import { removeFromCart } from '../actions/products';

class Cart extends React.Component {
  constructor(props) {
    super(props);

    this.redirectTo = this.redirectTo.bind(this);
  }

  redirectTo(path) {
    this.props.history.push(path);
  }

  render() {
    return (
      <div>
        <Navigation itemsInCart={this.props.itemsInCart} redirectTo={this.redirectTo} />

        <Grid>
          <Row>
            { this.props.cartItems.length === 0 ? (
              <h3>Nema proizvoda u košarici</h3>
          ) : (
            this.props.cartItems.map(item => (
              <CartCard
                key={item.id}
                item={item}
                removeFromCartAction={this.props.removeFromCart}
              />
           ))
          )
          }
          </Row>
        </Grid>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    cartItems: state.products.cartItems,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({
    removeFromCart,
  }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Cart);
