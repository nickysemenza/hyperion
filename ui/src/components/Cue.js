import React from "react";
import styled from "styled-components";
export const CueFrameWrapper = styled.div`
  display: flex;
`;

export const CueTable = styled.div`
  display: flex;
  border: 1px solid black;
`;
export const CueTableCol = styled.div`
  flex-direction: column;
`;

const TIME_SCALE = 0.3;

const CueFrameInner = styled.div.attrs({
  style: ({ duration }) => ({
    width: `${duration * TIME_SCALE}px`,
    minWidth: `${duration * TIME_SCALE}px`
  })
})`
  height: 50px;
  border: 1px solid #008aff;
  padding: 5px;
`;

export const CueFrame = ({ ...props }) => (
  <CueFrameInner {...props}>
    {" "}
    {props.duration} ms <br /> F:{props.frameId} | A:{props.actionId}
  </CueFrameInner>
);

const CueFrameWaitInner = styled.div.attrs({
  style: ({ duration }) => ({
    width: `${duration * TIME_SCALE}px`,
    minWidth: `${duration * TIME_SCALE}px`
  })
})`
  height: 50px;
  border: 1px solid purple;
  background-color: #f96f3a;
  padding: 5px;
`;
export const CueFrameWait = ({ ...props }) => (
  <CueFrameWaitInner {...props}> {props.duration} ms wait </CueFrameWaitInner>
);

const CueLabelInner = styled.div.attrs({
  style: ({ numActions }) => ({
    height: `${numActions * 50}px`
  })
})`
  // border: 1px solid purple;
  background-color: #20272b;
  width: 100px;
  color: white;
  display: flex;
  justify-content: center;
  flex-direction: column;
  text-align: center;
  order: 0;
`;

export const CueLabel = ({ ...props }) => (
  <CueLabelInner {...props}>
    Cue # {props.id} <br />
    {props.duration} ms
  </CueLabelInner>
);
