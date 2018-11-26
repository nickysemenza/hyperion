import { API_BASE_URL } from '../config';
import { RECEIVE_LIGHT_LIST, RECEIVE_STATE_LIST } from './lights';
import { RECEIVE_CUE_MASTER } from './cues';
export default function apiFetch(endpoint, options = {}) {
  options.headers = {
    // 'X-JWT': cookie.load('token')
  };
  return fetch(`${API_BASE_URL}/${endpoint}`, options);
}

const WS_TYPE_LIGHT_LIST = 'LIGHT_LIST';
const WS_TYPE_CUEMASTER = 'CUE_MASTER';
export const WS_META_OPEN = 'META_OPEN';
export const receiveSocketData = json => {
  // console.log("received socket data", json);
  let { data, type } = json;
  switch (type) {
    case WS_TYPE_LIGHT_LIST:
      return {
        type: RECEIVE_LIGHT_LIST,
        lights: data
      };
    case WS_TYPE_LIGHT_LIST:
      return {
        type: RECEIVE_LIGHT_LIST,
        lights: data
      };
    case 'LIGHT_STATES':
      return {
        type: RECEIVE_STATE_LIST,
        states: data
      };
    case WS_META_OPEN:
      return {
        type: WS_META_OPEN,
        open: data
      };
    default:
      return { type: null };
  }
};

export const UPDATE_WINDOW_DIMENSIONS = 'UPDATE_WINDOW_DIMENSIONS';
export const updateWindowDimensions = (width, height) => {
  return {
    type: UPDATE_WINDOW_DIMENSIONS,
    width,
    height
  };
};
