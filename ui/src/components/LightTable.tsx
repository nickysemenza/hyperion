import React from 'react';
import { ColorBox } from './LightState';
import { Table } from 'semantic-ui-react';
import { Light } from '../types';
type LightTableProps = {
  //TODO: make these more strongly typed
  lights: any;
  states: any;
};
const LightTable: React.SFC<LightTableProps> = ({ lights, states }) => {
  let tableRows = Object.keys(lights).map(k => {
    let eachLight = lights[k];
    let light = new Light(lights[k]);
    let eachState = states ? states[k] : null;
    let { name, state, ...meta } = eachLight;
    return (
      <Table.Row key={name}>
        <Table.Cell>{name}</Table.Cell>
        <Table.Cell>{light.getType()}</Table.Cell>
        <Table.Cell>
          <code>{JSON.stringify(meta, null, 2)}</code>
        </Table.Cell>
        <Table.Cell>
          <ColorBox state={eachState} />
        </Table.Cell>
        <Table.Cell>
          <code>{JSON.stringify(eachState)}</code>
        </Table.Cell>
      </Table.Row>
    );
  });

  return (
    <div>
      <Table singleLine>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>Name</Table.HeaderCell>
            <Table.HeaderCell>Type</Table.HeaderCell>
            <Table.HeaderCell>Meta</Table.HeaderCell>
            <Table.HeaderCell>Color</Table.HeaderCell>
            <Table.HeaderCell>State</Table.HeaderCell>
          </Table.Row>
        </Table.Header>

        <Table.Body>{tableRows}</Table.Body>
      </Table>
    </div>
  );
};
export default LightTable;
