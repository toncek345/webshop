import React from 'react';

import { Navbar, Nav, NavItem, Badge } from 'react-bootstrap';

import { Routes } from '../../constants';

const Navigation = ({
  redirectTo, itemsInCart,
}) => (
  <Navbar>
    <Navbar.Header>
      <Navbar.Brand>
        <a onClick={() => redirectTo(Routes.home.path)}>Generic webshop</a>
      </Navbar.Brand>
      <Navbar.Toggle />
    </Navbar.Header>
    <Navbar.Collapse>
      <Nav>
        <NavItem onClick={() => redirectTo(Routes.home.path)}>
        Products
        </NavItem>
      </Nav>
      <Nav pullRight>
        <NavItem onClick={() => redirectTo(Routes.cart.path)}>
          Cart <Badge>{itemsInCart}</Badge>
        </NavItem>
      </Nav>
    </Navbar.Collapse>
  </Navbar>
);

export default Navigation;
