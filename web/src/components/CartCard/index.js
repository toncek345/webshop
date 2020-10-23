import React from "react";

import { Panel, Button, Image, Col } from "react-bootstrap";

import { ServerIp } from "../../constants";

const CartCard = ({ item, removeFromCartAction }) => (
  <Panel>
    <Panel.Heading>{item.Name}</Panel.Heading>
    <Panel.Body>
      <Col xs={3} md={3}>
        <Image
          responsive
          src={`${ServerIp}/static/${item.images[0].path}`}
          style={{ height: 150 }}
        />
      </Col>
      <Col xs={6} md={6}>
        {item.name}
        <br />
        {item.price} kn
        <br />
        <Button onClick={() => removeFromCartAction(item)} bsStyle="danger">
          Makni iz košarice
        </Button>
      </Col>
    </Panel.Body>
  </Panel>
);

export default CartCard;
