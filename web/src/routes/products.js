import React from 'react';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { PageHeader, Grid, Row, Col, Modal, Image, Button, Carousel } from 'react-bootstrap';

import { getProducts, addToCart, removeFromCart } from '../actions/products';

import Navigation from '../components/Navigation';
import Card from '../components/Card';
import { ServerIp } from '../constants';

class Products extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedProduct: null,
      modal: false,
    };

    this.redirectTo = this.redirectTo.bind(this);
    this.handleProductSelect = this.handleProductSelect.bind(this);
    this.handleCloseModal = this.handleCloseModal.bind(this);
  }

  componentWillMount() {
    this.props.getProducts();
  }

  redirectTo(path) {
    this.props.history.push(path);
  }

  handleProductSelect(id) {
    const selectedProduct = this.props.products.find((product) => {
      if (product.Id === id) {
        return product;
      }
      return undefined; // linter satisfied
    });

    this.setState({
      selectedProduct,
      modal: true,
    });
  }

  handleCloseModal() {
    this.setState({
      selectedProduct: null,
      modal: false,
    });
  }

  isSelectedProductInCart() {
    if (this.props.cartItems.find(item => item.Id === this.state.selectedProduct.Id)) { return true; }

    return false;
  }

  render() {
    if (this.props.loading && !this.props.products) {
      return 'loading....';
    }

    if (this.props.error) {
      return this.props.error;
    }

    let modal = <Modal show={false} />;
    if (this.state.modal && this.state.selectedProduct) {
      modal = (
        <Modal show={this.state.modal} onHide={this.handleCloseModal} >
          <Modal.Header closeButton>
            <Modal.Title>{this.state.selectedProduct.Name}</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Carousel>

              {this.state.selectedProduct.Images.map(item => (
                <Carousel.Item key={item.Id}>
                  <Image src={`${ServerIp}/api/v1/static/${item.Name}`} />
                </Carousel.Item>
                  ))
              }

            </Carousel>
            <Grid>
              <Row>
                <h2>Cijena: {this.state.selectedProduct.Price / 100} kn  </h2>

              </Row>
              <Row>
                { this.isSelectedProductInCart() ? (
                  <Button
                    bsStyle="danger"
                    onClick={() => this.props.removeFromCart(this.state.selectedProduct)}
                  >Izbaci iz košarice
                  </Button>
                ) : (
                  <Button
                    bsStyle="primary"
                    onClick={() => this.props.addToCart(this.state.selectedProduct)}
                  >U košaricu
                  </Button>
                )}
              </Row>
            </Grid>
            <p />
            <hr />
            <div dangerouslySetInnerHTML={{ __html: this.state.selectedProduct.Description }} />

          </Modal.Body>
        </Modal>
      );
    }

    return (
      <div>
        <Navigation itemsInCart={this.props.itemsInCart} redirectTo={this.redirectTo} />

        <Grid>
          <Row>
            {this.props.products && this.props.products.map(item => (
              <Card
                key={item.Id}
                itemId={item.Id}
                imageUrl={item.Images && item.Images[0] && item.Images[0].Name}
                heading={item.Name}
                action={this.handleProductSelect}
                price={item.Price}
              />
              ))
            }
          </Row>
          {modal}
        </Grid>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    products: state.products.allProducts,
    loading: state.products.loading,
    error: state.products.error,

    cartItems: state.products.cartItems,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({
    getProducts,
    addToCart,
    removeFromCart,
  }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Products);
