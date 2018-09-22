import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { updateWindowDimensions } from '../actions';

class Helpers extends React.Component {
  componentDidMount = () => {
    this.updateWindowDimensions();
    window.addEventListener('resize', this.updateWindowDimensions);
  };

  componentWillUnmount = () => {
    window.removeEventListener('resize', this.updateWindowDimensions);
  };

  updateWindowDimensions = () => {
    this.props.updateWindowDimensions(window.innerWidth, window.innerHeight);
  };
  render = () => null;
}

function mapStateToProps(state) {
  return {};
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators({ updateWindowDimensions }, dispatch);
};

export default connect(mapStateToProps, mapDispatchToProps)(Helpers);
