import React from 'react';

import { Thumbnail, Button, Col } from 'react-bootstrap';

import { ServerIp } from '../../constants';

const Card = ({
  imageUrl, heading, action, itemId, price,
}) => (
  <Col xs={12} sm={4} md={3}>
    <Thumbnail
      src={`${ServerIp}/static/${imageUrl}`}
      href="#"
      onClick={() => action(itemId)}
      rounded
      responsive
    >
      <h3>{heading}</h3>
      <h3>{price && `${price / 100} kn`}</h3>
    </Thumbnail>
  </Col>
);

export default Card;
