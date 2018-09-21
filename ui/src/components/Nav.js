import React from 'react';
import { NavLink } from 'react-router-dom';
import { Container, Menu, Label } from 'semantic-ui-react';
import { connect } from 'react-redux';

const Nav = ({ user, ws_open }) => (
  <Menu fixed="top" inverted>
    <Container>
      <Menu.Item as={NavLink} to="/" exact header>
        hyperion
      </Menu.Item>
      <Menu.Item as={NavLink} to="/lights">
        Lights
      </Menu.Item>
      <Menu.Item as={NavLink} to="/cues">
        Cues
      </Menu.Item>
      <Menu.Item position="right">
        <Label color={ws_open ? 'green' : 'red'} horizontal>
          WebSocket
        </Label>
      </Menu.Item>
    </Container>
  </Menu>
);

function mapStateToProps(state) {
  return { user: state.user, ws_open: state.system.ws_open };
}
export default connect(mapStateToProps, {})(Nav);
