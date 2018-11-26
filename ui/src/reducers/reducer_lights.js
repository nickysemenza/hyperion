import { RECEIVE_LIGHT_LIST } from '../actions/lights';
import { RECEIVE_STATE_LIST } from '../actions/lights';
const INITIAL_STATE = {
  lights: { lights: {} },
  states: { states: {} }
};

export default function(state = INITIAL_STATE, action) {
  switch (action.type) {
    case RECEIVE_LIGHT_LIST:
      return {
        ...state,
        lights: action.lights
      };
    case RECEIVE_STATE_LIST:
      return {
        ...state,
        states: action.states
      };
    default:
      return state;
  }
}
