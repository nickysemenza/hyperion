import React, { Component } from 'react';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import LightTable from '../components/LightTable';
class LightList extends Component {
  render() {
    return <LightTable lights={this.props.lights} />;
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
