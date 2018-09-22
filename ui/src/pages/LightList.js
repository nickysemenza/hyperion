import React, { Component } from 'react';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Header } from 'semantic-ui-react';
import LightTable from '../components/LightTable';
class LightList extends Component {
  render() {
    return (
      <div>
        <Header
          as="h2"
          content="Lights"
          // subheader="View lights"
          icon="lightbulb blue"
          dividing
        />
        <LightTable lights={this.props.lights} />
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    lights: state.lights.lights
  };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators({}, dispatch);
};

export default connect(mapStateToProps, mapDispatchToProps)(LightList);
