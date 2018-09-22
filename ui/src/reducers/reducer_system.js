import { WS_META_OPEN, UPDATE_WINDOW_DIMENSIONS } from '../actions/index';
const INITIAL_STATE = {
  ws_open: false,
  windowDimensions: { width: 0, height: 0 }
};

export default function(state = INITIAL_STATE, action) {
  switch (action.type) {
    case WS_META_OPEN:
      return {
        ...state,
        ws_open: action.open
      };
    case UPDATE_WINDOW_DIMENSIONS:
      return {
        ...state,
        windowDimensions: {
          ...action
        }
      };
    default:
      return state;
  }
}
