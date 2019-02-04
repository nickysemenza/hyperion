import React, { Component } from 'react';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import LightTable from '../components/LightTable';
import { sendCommands } from '../actions/lights';
import { rgbToHex } from '../utils';
class Playground extends Component {
  state = {
    pressed: {
      83: true
    },
    colors: {
      r: 0,
      g: 0,
      b: 0
    }
  };
  handleKeyAction = (isDown, keyCode) => {
    console.log(keyCode, `is now ${isDown ? 'down' : 'up'}`);

    //mapping
    //81, 87, 69 QWE
    //65, 83, 68 ASD
    //90, 88, 67 ZXC

    let newVal;
    switch (keyCode) {
      case 81:
      case 87:
      case 69:
        newVal = 255;
        break;
      case 65:
      case 83:
      case 68:
        newVal = 128;
        break;
      default:
        newVal = 0;
        break;
    }
    let newChan;

    switch (keyCode) {
      case 65:
      case 81:
      case 90:
        newChan = 'r';
        break;
      case 87:
      case 83:
      case 88:
        newChan = 'g';
        break;
      case 69:
      case 68:
      case 67:
        newChan = 'b';
        break;
      default:
        newChan = '';
        break;
    }
    let { colors } = this.state;
    if (newChan !== '') {
      colors[newChan] = newVal;
      console.log(colors);
      let { r, g, b } = colors;
      // this.props.sendCommands([`set(hue2:${rgbToHex(r, g, b)}:0)`]);
      this.props.sendCommands([
        `set(hue2+par1:${rgbToHex(r, g, b)}+${rgbToHex(r, g, b)}:0+0)`
      ]);
      this.setState({ colors });
    }
  };
  handleKey = (isDown, e) => {
    let { keyCode } = e;
    // console.log("foo: ", isDown, keyCode);
    var curr = this.state.pressed[keyCode] || false;
    if (isDown !== curr) {
      var pressed = this.state.pressed;
      pressed[keyCode] = isDown;
      this.setState({ pressed });
      this.handleKeyAction(isDown, keyCode);
    }
  };

  render() {
    return (
      <div>
        <h2>lights</h2>
        <div
          tabIndex="0"
          onKeyDown={e => this.handleKey(true, e)}
          onKeyUp={e => this.handleKey(false, e)}
        >
          click me!
        </div>
        <LightTable lights={this.props.lights} states={this.props.states} />
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
  return bindActionCreators({ sendCommands }, dispatch);
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Playground);
