import React from 'react';
import { LinkContainer } from 'react-router-bootstrap';
import { Nav, Navbar, Badge } from 'react-bootstrap';
import { connect } from 'react-redux';

const MainNav = ({ user, ws_open }) => (
  <Navbar bg="dark" variant="dark">
    <LinkContainer to="/" exact>
      <Navbar.Brand>Hyperion</Navbar.Brand>
    </LinkContainer>
    <Nav className="mr-auto">
      <LinkContainer to="/lights">
        <Nav.Link>Lights</Nav.Link>
      </LinkContainer>
      <LinkContainer to="/cues">
        <Nav.Link>Cues</Nav.Link>
      </LinkContainer>
    </Nav>
    <Navbar.Text>
      <Badge variant={ws_open ? 'success' : 'error'}>websocket</Badge>
    </Navbar.Text>
  </Navbar>
);

function mapStateToProps(state) {
  return { user: state.user, ws_open: state.system.ws_open };
}
export default connect(
  mapStateToProps,
  {}
)(MainNav);
