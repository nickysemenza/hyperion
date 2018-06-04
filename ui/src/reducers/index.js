import { combineReducers } from 'redux';
import { reducer as toastrReducer } from 'react-redux-toastr';
import lights from './reducer_lights';
const rootReducer = combineReducers({
  lights,
  toastr: toastrReducer
});

export default rootReducer;
