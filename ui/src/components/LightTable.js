import React from 'react';
import { ColorBox } from '../components/LightState';
import { Table } from 'semantic-ui-react';
import { getLightType } from '../utils';
const LightTable = ({ lights, states }) => {
  let tableRows = Object.keys(lights).map(k => {
    let eachLight = lights[k];
    let eachState = states ? states[k] : null;
    let type = getLightType(eachLight);
    let { name, state, ...meta } = eachLight;
    return (
      <Table.Row key={name}>
        <Table.Cell>{name}</Table.Cell>
        <Table.Cell>{type}</Table.Cell>
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
