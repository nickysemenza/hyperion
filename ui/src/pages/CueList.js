import React, { Component } from "react";
import { connect } from "react-redux";
import { fetchCueMaster } from "../actions/cues";
import {
  CueTable,
  CueTableCol,
  CueLabel,
  CueFrame,
  CueFrameWait,
  CueFrameWrapper
} from "../components/Cue";
import { bindActionCreators } from "redux";
import { Header } from "semantic-ui-react";
class CueList extends Component {
  componentDidMount() {
    this.props.fetchCueMaster();
  }

  render() {
    let mainStack = this.props.stacks[0];

    if (!mainStack) return <div>loading...</div>;

    const bare = { length_ms: 0, items: [] };

    let all = {};
    mainStack.processed_cues.forEach(c => {
      let maxActions = 1;
      c.frames.forEach(
        f => (maxActions = Math.max(maxActions, f.actions.length))
      );
      all[c.id] = Array.apply(null, Array(maxActions)).map(x => bare);
    });

    mainStack.processed_cues.forEach(c =>
      c.frames.forEach(f => {
        f.actions.forEach((action, x) => {
          let tmp = {};
          Object.assign(tmp, all[c.id][x]);
          tmp["length_ms"] += action.action_duration_ms;
          tmp["items"] = tmp["items"].slice();
          tmp["items"].push(
            // "action" + action.id + "(" + action.action_duration_ms + "ms)"
            <CueFrame
              duration={action.action_duration_ms}
              frameId={f.id}
              actionId={action.id}
            />
          );
          all[c.id][x] = tmp;
        });
        //todo: add padding
        let targetLen = all[c.id][0]["length_ms"];
        all[c.id].forEach(
          x => (targetLen = Math.max(targetLen, x["length_ms"]))
        );
        all[c.id].forEach((item, x) => {
          let padding = targetLen - item["length_ms"];
          if (padding !== 0) {
            let tmp = {};
            Object.assign(tmp, all[c.id][x]);
            tmp["length_ms"] += padding;
            tmp["items"].push(
              // "padding:" + padding + "ms"

              <CueFrameWait duration={padding} />
            );
            all[c.id][x] = tmp;
          }
        });
      })
    );

    return (
      <div>
        <Header content={"cues"} />
        <CueTable>
          <CueTableCol>
            {mainStack.processed_cues.map(c => {
              let maxActions = 1;
              c.frames.forEach(
                f => (maxActions = Math.max(maxActions, f.actions.length))
              );
              return (
                <CueLabel
                  id={c.id}
                  numActions={maxActions}
                  duration={c.expected_duration_ms}
                />
              );
            })}
          </CueTableCol>
          <CueTableCol>
            {Object.keys(all).map(k => {
              let each = all[k];
              return each.map(item => (
                <CueFrameWrapper> {item.items}</CueFrameWrapper>
              ));
            })}
          </CueTableCol>
        </CueTable>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return { stacks: state.cues.cue_stacks };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators({ fetchCueMaster }, dispatch);
};

export default connect(mapStateToProps, mapDispatchToProps)(CueList);
