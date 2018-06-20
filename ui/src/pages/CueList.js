import React, { Component } from "react";
import { connect } from "react-redux";
import { fetchCueMaster } from "../actions/cues";
import Cue from "../components/Cue";
import { bindActionCreators } from "redux";
import { Header } from "semantic-ui-react";
class CueList extends Component {
  componentDidMount() {
    this.props.fetchCueMaster();
  }

  render() {
    let mainStack = this.props.stacks[0];

    if (!mainStack) return <div>loading...</div>;
    return (
      <div>
        <Header content={"cues"} />
        {this.props.stacks.length} stack(s)
        {mainStack.processed_cues.map(pc => <Cue key={pc.id} cue={pc} />)}
        {/* <pre>{JSON.stringify(mainStack, true, 2)}</pre> */}
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
