import React from 'react';
import { WS_URL } from '../config';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import Sockette from 'sockette';
import { receiveSocketData, WS_META_OPEN } from '../actions';

class WS extends React.Component {
  componentDidMount() {
    const ws = new Sockette(WS_URL, {
      // timeout: 5e3,
      onopen: e => {
        console.log('Connected!', e);
        this.props.receiveSocketData({ type: WS_META_OPEN, data: true });
        ws.send('hi');
      },
      onmessage: e => {
        try {
          let data = JSON.parse(e.data);
          this.props.receiveSocketData(data);
        } catch (error) {}
      },
      onreconnect: e => console.log('Reconnecting...', e),
      onmaximum: e => console.log('Stop Attempting!', e),
      onclose: e => {
        console.log('Closed!', e);
        this.props.receiveSocketData({ type: WS_META_OPEN, data: false });
      },
      onerror: e => console.log('Error:', e)
    });
  }
  render() {
    return null;
  }
}

function mapStateToProps(state) {
  return {};
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators({ receiveSocketData }, dispatch);
};

export default connect(mapStateToProps, mapDispatchToProps)(WS);
