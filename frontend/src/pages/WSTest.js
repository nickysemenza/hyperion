import React, { Component } from 'react';
import { API_BASE_URL } from '../config';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import Sockette from 'sockette';

class WSTest extends Component {
  componentDidMount() {
    const ws = new Sockette('ws://localhost:8080/ws', {
      timeout: 5e3,
      maxAttempts: 10,
      onopen: e => {
        console.log('Connected!', e);
        ws.send('hi');
      },
      onmessage: e => console.log('Received:', e),
      onreconnect: e => console.log('Reconnecting...', e),
      onmaximum: e => console.log('Stop Attempting!', e),
      onclose: e => console.log('Closed!', e),
      onerror: e => console.log('Error:', e)
    });
  }
  render() {
    return <div>hello</div>;
  }
}

function mapStateToProps(state) {
  return {
    config: state.lights.lights
  };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators({}, dispatch);
};

export default connect(mapStateToProps, mapDispatchToProps)(WSTest);
