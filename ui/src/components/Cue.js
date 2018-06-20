import React from "react";
import { Table } from "semantic-ui-react";

const Cue = ({ cue }) => {
  let { frames, ...meta } = cue;
  return (
    <div>
      CUE: # {cue.id}
      <br />
      number of frames: {cue.frames.length} <br />
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
