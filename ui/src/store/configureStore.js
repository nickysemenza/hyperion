import { createStore, applyMiddleware, compose } from 'redux';
import thunk from 'redux-thunk';
// import { createLogger } from 'redux-logger';
import rootReducer from '../reducers';

const configureStoreDev = () => {
  const composeEnhancers =
    typeof window === 'object' && window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__
      ? window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__({
          // Specify extension’s options like name, actionsBlacklist, actionsCreators, serialize...
        })
      : compose;

  const enhancer = composeEnhancers(
    applyMiddleware(
      thunk
      // createLogger()
    )
  );
  const store = createStore(rootReducer, {}, enhancer);

  if (module.hot) {
    // Enable Webpack hot module replacement for reducers
    module.hot.accept('../reducers', () => {
      const nextRootReducer = require('../reducers').default;
      store.replaceReducer(nextRootReducer);
    });
  }

  return store;
};
let store = {};
if (process.env.NODE_ENV === 'production') {
  store = createStore(rootReducer, {}, applyMiddleware(thunk));
  // module.exports = require('./configureStore.prod.js');
} else {
  store = configureStoreDev();
}
export default store;
