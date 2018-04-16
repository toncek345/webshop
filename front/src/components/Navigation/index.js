import React from 'react';

import { Navbar, Nav, NavItem, Badge } from 'react-bootstrap';

import { Routes } from '../../constants';

const Navigation = ({
  redirectTo, itemsInCart,
}) => (
  <Navbar>
    <Navbar.Header>
      <Navbar.Brand>
        <a href="#" onClick={() => redirectTo(Routes.home.path)}>Genericki webshop</a>
      </Navbar.Brand>
      <Navbar.Toggle />
    </Navbar.Header>
    <Navbar.Collapse>
      <Nav>
        <NavItem href="#" onClick={() => redirectTo(Routes.home.path)}>
        News
        </NavItem>
        <NavItem href="#" onClick={() => redirectTo(Routes.products.path)}>
        Products
        </NavItem>
      </Nav>
      <Nav pullRight>
        <NavItem href="#" onClick={() => redirectTo(Routes.cart.path)}>
          Cart <Badge>{itemsInCart}</Badge>
        </NavItem>
      </Nav>
    </Navbar.Collapse>
  </Navbar>
);

export default Navigation;
