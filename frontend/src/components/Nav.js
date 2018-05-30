import React from 'react';
import { NavLink } from 'react-router-dom';
import { Container, Menu } from 'semantic-ui-react';
import { connect } from 'react-redux';

const Nav = ({ user }) => (
  <Menu fixed="top" inverted>
    <Container>
      <Menu.Item as={NavLink} to="/" exact header>
        hyperion
      </Menu.Item>
      <Menu.Item as={NavLink} to="/lights">
        Lights
      </Menu.Item>
    </Container>
  </Menu>
);

function mapStateToProps(state) {
  return { user: state.user };
}
export default connect(mapStateToProps, {})(Nav);
