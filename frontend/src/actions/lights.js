import apiFetch from './index';

export const RECEIVE_LIGHT_LIST = 'RECEIVE_LIGHT_LIST';

export function fetchLightList() {
  return dispatch => {
    return apiFetch('lights')
      .then(response => response.json())
      .then(json => dispatch(receiveLightList(json)));
  };
}
function receiveLightList(lights) {
  return {
    type: RECEIVE_LIGHT_LIST,
    lights
  };
}
