import {RECEIVE_LIGHT_LIST} from '../actions/lights'
const INITIAL_STATE = {
    lights: {lights: {}}
}

export default function (state = INITIAL_STATE, action) {
    switch(action.type) {
        case RECEIVE_LIGHT_LIST:
        return {
            ...state,
            lights: action.lights
        }
        default: 
        return state;
    }
}