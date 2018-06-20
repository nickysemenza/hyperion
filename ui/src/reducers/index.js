import { combineReducers } from "redux";
import { reducer as toastrReducer } from "react-redux-toastr";
import lights from "./reducer_lights";
import cues from "./reducer_cues";
const rootReducer = combineReducers({
  lights,
  cues,
  toastr: toastrReducer
});

export default rootReducer;
