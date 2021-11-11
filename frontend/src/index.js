import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';

const title = 'Starboard Reports';

window._env_ = {
  REPORT_URL: window.location.href,
}

ReactDOM.render(
  <App title={title} />,
  document.getElementById('app')
);

module.hot.accept();
