import React from 'react';
import styled from 'styled-components';
import { rgbToHex, isRGBDark } from '../utils';

const COLOR_PROCESSED = '#2c3e50';
const COLOR_ENQUEUED = '#2980b9';
const COLOR_ACTIVE = '#2ecc71';
const COLOR_WAIT = '#e74c3c';
export const CueFrameWrapper = styled.div`
  display: flex;
`;

export const CueTable = styled.div`
  display: flex;
  border: 1px solid black;
  overflow-x: auto;
`;
export const CueTableCol = styled.div`
  flex-direction: column;
`;

const CueFrameInner = styled.div.attrs({
  style: ({ width, hex, contrastingTextColor }) => ({
    width: `${width}px`,
    minWidth: `${width}px`,
    backgroundColor: hex,
    color: contrastingTextColor
  })
})`
  height: 33px;
  border: 1px solid black;
  padding: 3px;
  font-size: 12px;
`;

export const CueFrame = ({ ...props }) => {
  let { rgb } = props.action.new_state;
  let { r, g, b } = rgb;

  props.hex = rgbToHex(r, g, b);
  props.contrastingTextColor = isRGBDark(r, g, b) ? '#292937' : '#F7F7F7';
  return (
    <CueFrameInner {...props}>
      {props.duration} ms {' | '}
      {props.debug ? ` (F:${props.frameId} | A:${props.actionId}) | ` : null}
      <b>{props.action.light_name}</b>
    </CueFrameInner>
  );
};

const CueFrameWaitInner = styled.div.attrs({
  style: ({ width }) => ({
    width: `${width}px`,
    minWidth: `${width}px`
  })
})`
  height: 33px;
  border: 1px solid black;
  background-color: ${COLOR_WAIT};
  padding: 5px;
`;
export const CueFrameWait = ({ ...props }) => (
  <CueFrameWaitInner {...props}> {props.duration} ms wait </CueFrameWaitInner>
);

const CueLabelInner = styled.div.attrs({
  style: ({ numActions, status }) => {
    let statusColor = COLOR_ENQUEUED;
    if (status === 'active') statusColor = COLOR_ACTIVE;
    else if (status === 'processed') statusColor = COLOR_PROCESSED;
    return {
      height: `${numActions * 33 + 16}px`,
      backgroundColor: statusColor
    };
  }
})`
  width: 100px;
  border: 1px solid black;
  color: white;
  display: flex;
  justify-content: center;
  flex-direction: column;
  text-align: center;
  order: 0;
`;

export const CueLabel = ({ ...props }) => {
  let { id, duration, duration_drift_ms, cue, debug } = props;
  return (
    <CueLabelInner {...props}>
      # {id} <br />
      {`${duration} ms`}{' '}
      {cue.status === 'active'
        ? `${((cue.elapsed_ms / cue.expected_duration_ms) * 100).toFixed(1)} %`
        : null}
      <i>
        {(debug && duration_drift_ms && `(+${duration_drift_ms} ms)`) || null}
      </i>
    </CueLabelInner>
  );
};

const ProgressInner = styled.div.attrs({
  style: ({ width, color }) => ({
    width: `${width}px`,
    minWidth: `${width}px`,
    backgroundColor: color
  })
})`
  height: 16px;
  border: 1px solid black;
`;
export const Progress = ({ ...props }) => {
  let { cue } = props;
  let { elapsed_ms, expected_duration_ms, status } = cue;
  let duration = 0;
  let color = COLOR_PROCESSED;
  if (status === 'active') {
    duration = Math.min(elapsed_ms, expected_duration_ms);
    color = COLOR_ACTIVE;
  }
  if (status === 'processed') duration = expected_duration_ms;

  return (
    <ProgressInner
      width={props.scaleDurationToUIWidth(duration)}
      color={color}
    />
  );
};
