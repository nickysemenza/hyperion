import { combineReducers } from 'redux';
import { reducer as toastrReducer } from 'react-redux-toastr';
import lights from './reducer_lights';
import cues from './reducer_cues';
import system from './reducer_system';
const rootReducer = combineReducers({
  lights,
  cues,
  system,
  toastr: toastrReducer
});

export default rootReducer;
