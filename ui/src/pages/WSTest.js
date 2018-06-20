import React, { Component } from "react";
import { WS_URL } from "../config";
import moment from "moment";
import { bindActionCreators } from "redux";
import { connect } from "react-redux";
import Sockette from "sockette";
import LightTable from "../components/LightTable";
import { Label } from "semantic-ui-react";

class WSTest extends Component {
  state = {
    lights: {},
    curTime: null,
    wsOpen: false
  };
  componentDidMount() {
    setInterval(() => {
      this.setState({
        //   curTime : moment().valueOf()
        curTime: moment().format("HH:mm:ss:SS (x)")
      });
    }, 10);

    const ws = new Sockette(WS_URL, {
      // timeout: 5e3,
      onopen: e => {
        console.log("Connected!", e);
        this.setState({ wsOpen: true });
        ws.send("hi");
      },
      onmessage: e => {
        // console.log('Received:', e);
        try {
          let lights = JSON.parse(e.data);
          this.setState({ lights });
        } catch (error) {}
      },
      onreconnect: e => console.log("Reconnecting...", e),
      // onmaximum: e => console.log('Stop Attempting!', e),
      onclose: e => {
        console.log("Closed!", e);
        this.setState({ wsOpen: false });
      },
      onerror: e => console.log("Error:", e)
    });
  }
  render() {
    let { lights, wsOpen } = this.state;
    return (
      <div>
        <Label color={wsOpen ? "green" : "red"} horizontal>
          WebSocket
        </Label>
        <br />
        <LightTable lights={lights} />
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
