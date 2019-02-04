import React, { Component } from 'react';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import LightTable from '../components/LightTable';
class LightList extends Component {
  render() {
    return (
      <div>
        <h2>lights</h2>
        <LightTable
          foo="bar"
          lights={this.props.lights}
          states={this.props.states}
        />
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    lights: state.lights.lights,
    states: state.lights.states
  };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators({}, dispatch);
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(LightList);
