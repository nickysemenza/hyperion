import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';

import Home from './pages/static/Home';
import LightList from './pages/LightList';
import CueList from './pages/CueList';
import Nav from './components/Nav';
import WS from './components/WS';
import Helpers from './components/Helpers';
import { Container } from 'react-bootstrap';
import ReduxToastr from 'react-redux-toastr';
import Playground from './pages/Playground';

const App = () => (
  <Router>
    <div>
      <ReduxToastr
        timeOut={4000}
        newestOnTop={false}
        preventDuplicates
        position="top-left"
        transitionIn="fadeIn"
        transitionOut="fadeOut"
        progressBar
      />
      <Nav />
      <Container
        fluid
        style={{ marginTop: '7em', width: '95%', minHeight: '100vh' }}
      >
        <Switch>
          <Route exact path="/" component={Home} />
          <Route path="/playground" component={Playground} />
          <Route path="/lights" component={LightList} />
          <Route path="/cues" component={CueList} />
        </Switch>
        <WS />
        <Helpers />
      </Container>
    </div>
  </Router>
);
export default App;
