import { RECEIVE_CUE_MASTER } from '../actions/cues';
const INITIAL_STATE = {
  cue_stacks: []
};

export default function(state = INITIAL_STATE, action) {
  switch (action.type) {
    case RECEIVE_CUE_MASTER:
      return {
        ...state,
        cue_stacks: action.cuemaster.cue_stacks
      };
    default:
      return state;
  }
}
