import React from 'react';
import ReactDOM from 'react-dom';
import './assets/index.css';
import App from './App';
import registerServiceWorker from './registerServiceWorker';
import configureStore from './store/configureStore';
import { Provider } from 'react-redux';

import 'react-redux-toastr/lib/css/react-redux-toastr.min.css';
import './semantic/dist/semantic.min.css';

const store = configureStore();

ReactDOM.render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('root')
);
registerServiceWorker();
