import React, { Component } from 'react';

import { Table } from 'semantic-ui-react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { fetchLightList } from '../actions/lights';

class LightList extends Component {
  componentDidMount() {
    this.props.fetchLightList();
  }
  render() {
    let { config } = this.props;
    let { lights } = config;

    let data = [];
    const types = ['hue', 'dmx', 'generic'];
    for (let x in types) {
      let type = types[x];
      let lightsOfType = lights[type];

      for (let y in lightsOfType) {
        let eachLight = lightsOfType[y];
        let { name, ...meta } = eachLight;
        data.push({ name, meta, type, key: name });
      }
    }
    const tableRows = data.map(x => (
      <Table.Row key={x.name}>
        <Table.Cell>{x.name}</Table.Cell>
        <Table.Cell>{x.type}</Table.Cell>
        <Table.Cell>
          <pre>{JSON.stringify(x.meta, null, 2)}</pre>
        </Table.Cell>
      </Table.Row>
    ));

    return (
      <div>
        <Table celled>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Name</Table.HeaderCell>
              <Table.HeaderCell>Type</Table.HeaderCell>
              <Table.HeaderCell>Meta</Table.HeaderCell>
            </Table.Row>
          </Table.Header>

          <Table.Body>{tableRows}</Table.Body>
        </Table>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    config: state.lights.lights
  };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators(
    {
      fetchLightList
    },
    dispatch
  );
};

export default connect(mapStateToProps, mapDispatchToProps)(LightList);
