import React from 'react';

const componentToHex = c => {
  var hex = c.toString(16);
  return hex.length == 1 ? '0' + hex : hex;
};

const rgbToHex = (r, g, b) => {
  return '#' + componentToHex(r) + componentToHex(g) + componentToHex(b);
};

const LightState = ({ s, name }) => {
  let c = s['rgb'];
  return (
    <div>
    {name}
    <div style={{ backgroundColor: rgbToHex(c.r, c.g, c.b), padding: "20px", width:"50px"}}/>
      {/* <pre>{JSON.stringify(s, null, 2)}</pre> */}
    </div>
  );
};
export default LightState;
