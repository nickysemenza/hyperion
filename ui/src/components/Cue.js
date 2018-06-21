import React from "react";
import { Table } from "semantic-ui-react";
import styled from "styled-components";
const Cue = ({ cue }) => {
  let { frames, ...meta } = cue;

  const CueFrameWrapper = styled.div`
    display: flex;
  `;

  const CueTable = styled.div`
    display: flex;
    border: 1px solid black;
  `;
  const Col = styled.div`
    flex-direction: column;
  `;

  const CueFrameInner = styled.div.attrs({
    style: ({ duration }) => ({
      width: `${duration}px`,
      minWidth: `${duration}px`
    })
  })`
    height: 50px;
    border: 1px solid #008aff;
  `;

  const CueFrame = ({ ...props }) => (
    <CueFrameInner {...props}> {props.duration} ms</CueFrameInner>
  );

  const CueFrameWaitInner = styled.div.attrs({
    style: ({ duration }) => ({
      width: `${duration}px`,
      minWidth: `${duration}px`
    })
  })`
    height: 50px;
    border: 1px solid purple;
    background-color: #f96f3a;
  `;
  const CueFrameWait = ({ ...props }) => (
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

  const CueLabel = ({ ...props }) => (
    <CueLabelInner {...props}>
      Cue # {props.id} <br />XXX ms
    </CueLabelInner>
  );

  return (
    <div>
      CUE: # {cue.id}
      <br />
      number of frames: {cue.frames.length} <br />
      <CueTable>
        <Col>
          <CueLabel numActions={2} id={1} />
          <CueLabel numActions={1} id={2} />
        </Col>
        <Col>
          <CueFrameWrapper>
            <CueFrame duration={100} />
            <CueFrame duration={150} />
          </CueFrameWrapper>
          <CueFrameWrapper>
            <CueFrame duration={100} />
            <CueFrameWait duration={50} />
            <CueFrame duration={100} />
            <CueFrame duration={50} />
          </CueFrameWrapper>
          <CueFrameWrapper>
            <CueFrameWait duration={275} />
            <CueFrame duration={25} />
          </CueFrameWrapper>
        </Col>
      </CueTable>
      <hr />
      <Table celled>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>ID</Table.HeaderCell>
            <Table.HeaderCell>Frames</Table.HeaderCell>
            {/* <Table.HeaderCell>Meta</Table.HeaderCell> */}
            {/* <Table.HeaderCell>State</Table.HeaderCell> */}
          </Table.Row>
        </Table.Header>

        <Table.Body>
          {frames.map(cf => (
            <Table.Row key={cf.id}>
              <Table.Cell>{cf.id}</Table.Cell>
              {/* <Table.Cell>{type}</Table.Cell> */}
              <Table.Cell>
                <pre>{JSON.stringify(cf, null, 2)}</pre>
              </Table.Cell>
              {/* <Table.Cell> */}
              {/* <LightState s={state} name={name} /> */}
              {/* </Table.Cell> */}
            </Table.Row>
          ))}
        </Table.Body>
      </Table>
      <pre>{JSON.stringify(meta, true, 2)}</pre>
    </div>
  );
};

export default Cue;
