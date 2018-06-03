import React, { Component } from 'react';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { fetchLightList } from '../actions/lights';

import LightTable from '../components/LightTable';
class LightList extends Component {
  componentDidMount() {
    this.props.fetchLightList();
  }
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
  return bindActionCreators(
    {
      fetchLightList
    },
    dispatch
  );
};

export default connect(mapStateToProps, mapDispatchToProps)(LightList);
