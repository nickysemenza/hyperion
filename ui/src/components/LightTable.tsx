import React from 'react';
import { ColorBox } from './LightState';
import { Table } from 'react-bootstrap';
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
      <tr key={name}>
        <td>{name}</td>
        <td>{light.getType()}</td>
        <td>
          <code>{JSON.stringify(meta, null, 2)}</code>
        </td>
        <td>
          <ColorBox state={eachState} />
        </td>
        <td>
          <code>{JSON.stringify(eachState)}</code>
        </td>
      </tr>
    );
  });

  return (
    <div>
      <Table striped bordered hover>
        <thead>
          <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Meta</th>
            <th>Color</th>
            <th>State</th>
          </tr>
        </thead>

        <tbody>{tableRows}</tbody>
      </Table>
    </div>
  );
};
export default LightTable;
