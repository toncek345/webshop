import React from 'react';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { PageHeader, Grid, Row, Modal, Image } from 'react-bootstrap';

import { getNews } from '../actions/news';
import Navigation from '../components/Navigation';
import Card from '../components/Card';
import { ServerIp } from '../constants';

// home component shows news & navigation + clicking on any of news opens modal
class Home extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      modal: false,
      selectedNews: null,
    };

    this.redirectTo = this.redirectTo.bind(this);
    this.handleNewsClick = this.handleNewsClick.bind(this);
    this.handleCloseModal = this.handleCloseModal.bind(this);
  }

  componentWillMount() {
    this.props.getNews();
  }

  redirectTo(path) {
    this.props.history.push(path);
  }

  handleNewsClick(id) {
    const selectedNews = this.props.news.find((element) => {
      if (element.Id === id) {
        return element;
      }
      return undefined; // linter satisfied
    });

    this.setState({
      selectedNews,
      modal: true,
    });
  }

  handleCloseModal() {
    this.setState({
      selectedNews: null,
      modal: false,
    });
  }

  render() {
    if (this.props.error) {
      return this.props.error;
    }

    if (this.props.loading && !this.props.news) {
      return 'loading...';
    }

    let modal = <Modal show={false} />;
    if (this.state.modal && this.state.selectedNews) {
      modal = (
        <Modal show={this.state.modal} onHide={this.handleCloseModal} >
          <Modal.Header closeButton>
            <Modal.Title>{this.state.selectedNews.Header}</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Image src={`${ServerIp}/static/${this.state.selectedNews.ImagePath}`} responsive />
            <hr />
            <p>{this.state.selectedNews.Text}</p>
          </Modal.Body>
        </Modal>
      );
    }

    return (
      <div>
        <Navigation itemsInCart={this.props.itemsInCart} redirectTo={this.redirectTo} />

        <Grid>
          <Row>
            {this.props.news && this.props.news.map(item => (
              <Card
                key={item.Id}
                itemId={item.Id}
                imageUrl={item.ImagePath}
                heading={item.Header}
                action={this.handleNewsClick}
              />))
            }
          </Row>
        </Grid>

        {modal}
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    news: state.news.allNews,
    loading: state.news.loading,
    error: state.news.error,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({
    getNews,
  }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Home);
