import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';

const title = 'My Minimal React Webpack Babel Setup';

window._env_ = {
  SERVER_URL: "http://9.30.189.42:8889/",
  REPORT_URL: "http://9.30.189.42:8888/",
}

ReactDOM.render(
  <App title={title} />,
  document.getElementById('app')
);

module.hot.accept();
