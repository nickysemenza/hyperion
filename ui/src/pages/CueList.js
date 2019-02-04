import React, { Component } from 'react';
import { connect } from 'react-redux';
import { fetchCueMaster } from '../actions/cues';
import {
  CueTable,
  CueTableCol,
  CueLabel,
  CueFrame,
  CueFrameWait,
  CueFrameWrapper,
  Progress
} from '../components/Cue';
import { bindActionCreators } from 'redux';
import { Button, ButtonGroup } from 'react-bootstrap';
class cueList extends Component {
  state = {
    scale: 0.8,
    debug: false,
    autoZoom: true,
    longestCueDuration: 0.0
  };
  componentDidMount() {
    this.props.fetchCueMaster();
  }
  scaleDurationToUIWidth = duration => duration * this.state.scale;

  changeZoom = direction => {
    //100px for the left column, 40px padding on either size, 5px buffer
    let cueDisplayWidth = this.props.windowDimensions.width - 185;
    if (direction === 'auto') {
      this.setState({
        scale: cueDisplayWidth / this.state.longestCueDuration,
        autoZoom: true
      });
      return;
    }
    this.setState({
      scale: this.state.scale * (direction === 'in' ? 1.2 : 0.8),
      autoZoom: false
    });
  };
  toggleDebug = () => {
    this.setState({ debug: !this.state.debug });
  };
  setLongestCueDuration = duration => {
    if (duration > this.state.longestCueDuration)
      this.setState({ longestCueDuration: duration });
  };

  componentDidUpdate = (prevProps, prevState, snapshot) => {
    let oldStack = this.getMainStack(prevProps.stacks);
    let newStack = this.getMainStack(this.props.stacks);
    if (oldStack && newStack && oldStack !== newStack && this.state.autoZoom)
      this.changeZoom('auto');
  };

  //todo: move to selector
  getMainStack = stacks => stacks && stacks[0] && stacks[0].processed_cues;

  render() {
    let mainStack = this.props.stacks[0];

    if (!mainStack) return <div>loading...</div>;

    const bare = { length_ms: 0, items: [] };

    let all = {};
    let cuesById = {};
    if (!mainStack.processed_cues) mainStack.processed_cues = [];

    let cueList = mainStack.processed_cues.concat(
      mainStack.cues || [],
      mainStack.active_cue || []
    );

    cueList.sort((a, b) => a.id - b.id);

    cueList.forEach(c => {
      let maxActions = 1;
      c.frames.forEach(
        f => (maxActions = Math.max(maxActions, f.actions.length))
      );
      all[c.id] = Array.apply(null, Array(maxActions)).map(x => bare);
      cuesById[c.id] = c;
    });

    cueList.forEach(c =>
      c.frames.forEach((f, z) => {
        f.actions.forEach((action, x) => {
          let tmp = {};
          Object.assign(tmp, all[c.id][x]);
          tmp['length_ms'] += action.action_duration_ms;
          tmp['items'] = tmp['items'].slice();
          tmp['items'].push(
            <CueFrame
              action={action}
              width={this.scaleDurationToUIWidth(action.action_duration_ms)}
              key={c.id + '-' + z + '-' + x}
              duration={action.action_duration_ms}
              frameId={f.id}
              actionId={action.id}
              debug={this.state.debug}
            />
          );
          all[c.id][x] = tmp;
        });
        //todo: add padding
        //figure out how long the whole box needs to be
        let targetLen = all[c.id][0]['length_ms'];
        this.setLongestCueDuration(targetLen);
        all[c.id].forEach(
          x => (targetLen = Math.max(targetLen, x['length_ms']))
        );
        //add padding where necessary
        all[c.id].forEach((item, x) => {
          let padding = targetLen - item['length_ms'];
          if (padding !== 0) {
            let tmp = {};
            Object.assign(tmp, all[c.id][x]);
            tmp['length_ms'] += padding;
            tmp['items'].push(
              <CueFrameWait
                key={x}
                duration={padding}
                width={this.scaleDurationToUIWidth(padding)}
              />
            );
            all[c.id][x] = tmp;
          }
        });
      })
    );

    return (
      <div>
        <h2>cues</h2>
        {/* zoom buttons */}
        <ButtonGroup>
          <Button onClick={() => this.changeZoom('out')}>zoom out</Button>
          <Button onClick={() => this.changeZoom('in')}>zoom in</Button>
          <Button
            active={this.state.autoZoom}
            onClick={() => this.changeZoom('auto')}
          >
            auto-zoom
          </Button>
          <Button active={this.state.debug} onClick={this.toggleDebug}>
            debug
          </Button>
        </ButtonGroup>
        <hr />
        <CueTable>
          <CueTableCol>
            {cueList.map(c => {
              let maxActions = 1;
              c.frames.forEach(
                f => (maxActions = Math.max(maxActions, f.actions.length))
              );
              return (
                <CueLabel
                  id={c.id}
                  key={c.id}
                  numActions={maxActions}
                  status={c.status}
                  cue={c}
                  duration={c.expected_duration_ms}
                  duration_drift_ms={c.duration_drift_ms}
                  debug={this.state.debug}
                />
              );
            })}
          </CueTableCol>
          <CueTableCol>
            {Object.keys(all).map(k => {
              return (
                <div key={k + '-23'}>
                  <CueFrameWrapper key={k + '-2'}>
                    <Progress
                      cue={cuesById[parseInt(k, 10)]}
                      scaleDurationToUIWidth={this.scaleDurationToUIWidth}
                    />
                  </CueFrameWrapper>
                  {all[k].map((item, x) => (
                    <CueFrameWrapper key={x}> {item.items}</CueFrameWrapper>
                  ))}
                </div>
              );
            })}
          </CueTableCol>
        </CueTable>
      </div>
    );
  }
}

function mapStateToProps(state) {
  let { windowDimensions } = state.system;
  return { stacks: state.cues.cue_stacks, windowDimensions };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators({ fetchCueMaster }, dispatch);
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(cueList);
