import { API_BASE_URL } from "../config";
import { RECEIVE_LIGHT_LIST } from "./lights";
import { RECEIVE_CUE_MASTER } from "./cues";
export default function apiFetch(endpoint, options = {}) {
  options.headers = {
    // 'X-JWT': cookie.load('token')
  };
  return fetch(`${API_BASE_URL}/${endpoint}`, options);
}

const WS_TYPE_LIGHT_LIST = "LIGHT_LIST";
const WS_TYPE_CUEMASTER = "CUE_MASTER";
export const receiveSocketData = json => {
  // console.log("received socket data", json);
  let { data, type } = json;
  switch (type) {
    case WS_TYPE_LIGHT_LIST:
      return {
        type: RECEIVE_LIGHT_LIST,
        lights: data
      };
    case WS_TYPE_CUEMASTER:
      return {
        type: RECEIVE_CUE_MASTER,
        cuemaster: data
      };
    default:
      return { type: null };
  }
};
