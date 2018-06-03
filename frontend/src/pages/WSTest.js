import React, { Component } from 'react';
import { WS_URL } from '../config';
import moment from 'moment';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import LightState from '../components/LightState';
import Sockette from 'sockette';

class WSTest extends Component {
  state = {
    lights: {},
    curTime: null
  };
  componentDidMount() {
    setInterval(() => {
      this.setState({
        //   curTime : moment().valueOf()
        curTime: moment().format('HH:mm:ss:SS (x)')
      });
    }, 10);

    const ws = new Sockette(WS_URL, {
      timeout: 5e3,
      maxAttempts: 10,
      onopen: e => {
        console.log('Connected!', e);
        ws.send('hi');
      },
      onmessage: e => {
        // console.log('Received:', e);
        try {
          let lights = JSON.parse(e.data);
          this.setState({ lights });
        } catch (error) {}
      },
      onreconnect: e => console.log('Reconnecting...', e),
      onmaximum: e => console.log('Stop Attempting!', e),
      onclose: e => console.log('Closed!', e),
      onerror: e => console.log('Error:', e)
    });
  }
  render() {
    let { lights } = this.state;
    let lightDetail = Object.keys(lights).map((k, i) => {
      let eachLight = lights[k];
      return (
        <LightState key={i} s={eachLight['state']} name={eachLight['name']} />
      );
    });
    return (
      <div>
        hello
        {lightDetail}
        <pre>{JSON.stringify(this.state.curTime, null, 2)}</pre>
      </div>
    );
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
