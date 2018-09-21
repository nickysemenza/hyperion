import { WS_META_OPEN } from '../actions/index';
const INITIAL_STATE = {
  ws_open: false
};

export default function(state = INITIAL_STATE, action) {
  switch (action.type) {
    case WS_META_OPEN:
      return {
        ...state,
        ws_open: action.open
      };
    default:
      return state;
  }
}
