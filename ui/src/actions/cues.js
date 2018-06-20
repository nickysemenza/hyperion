import apiFetch from "./index";

export const RECEIVE_CUE_MASTER = "RECEIVE_CUE_MASTER";

export function fetchCueMaster() {
  return dispatch => {
    return apiFetch("cuemaster")
      .then(response => response.json())
      .then(json => dispatch(receiveCueMaster(json)));
  };
}
function receiveCueMaster(cuemaster) {
  return {
    type: RECEIVE_CUE_MASTER,
    cuemaster
  };
}
