import React, { Component } from 'react';
import './App.css';
import { Table} from 'antd';

class App extends Component {
  state = {
    a: {lights: {}}
  };
  componentDidMount() {
    fetch("http://localhost:8080/lights")
    .then(function(response) {
      return response.json();
    })
    .then(d=>{
      console.log(d);
      this.setState({a: d})
    })
  }
  render() {

    let data2 = []
    const types = ["hue","dmx","generic"];
    for(let x in types) {
      let type = types[x];
      let lightsOfType = this.state.a.lights[type]

      for(let y in lightsOfType) {
        let eachLight = lightsOfType[y]
        let {name, ...meta} = eachLight
        data2.push({name, meta, type, key: name})
      }
    }
    
    const columns = [{
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },{
      title: 'Type',
      dataIndex: 'type',
      key: 'type',
    }, {
      title: 'Meta',
      key: 'meta',
      render: (text, record) => (
        <span>
         <pre>{JSON.stringify(text,null,2)}</pre> 
        </span>
      ), 
    }];
    
    return (
      <div style={{ margin: 24 }}>
    <div style={{ marginBottom: 24 }}>
    <h1>Lights</h1>
    <Table dataSource={data2} columns={columns} />
    {/* <pre>{JSON.stringify(this.state,null,2)}</pre> */}
    </div>
  </div>
    );
  }
}

export default App;
