import React from 'react';
import LightState from '../components/LightState';
import { Table } from 'semantic-ui-react';
import { getLightType } from '../utils';
const LightTable= ({lights}) => {
    console.log("hi");
    let tableRows = Object.keys(lights).map(k => {
        let eachLight = lights[k];
        let type = getLightType(eachLight);
        let { name, state, ...meta } = eachLight;
        return (
          <Table.Row key={name}>
            <Table.Cell>{name}</Table.Cell>
            <Table.Cell>{type}</Table.Cell>
            <Table.Cell>
              <pre>{JSON.stringify(meta, null, 2)}</pre>
            </Table.Cell>
            <Table.Cell>
              <LightState s={state} name={name} />
            </Table.Cell>
          </Table.Row>
        )
      });
  
      return (
        <div>
          <Table celled>
            <Table.Header>
              <Table.Row>
                <Table.HeaderCell>Name</Table.HeaderCell>
                <Table.HeaderCell>Type</Table.HeaderCell>
                <Table.HeaderCell>Meta</Table.HeaderCell>
                <Table.HeaderCell>State</Table.HeaderCell>
              </Table.Row>
            </Table.Header>
  
            <Table.Body>{tableRows}</Table.Body>
          </Table>
        </div>
      );
};
export default LightTable;
