import React from 'react';
import styled from 'styled-components';
import { rgbToHex } from '../utils';

const SampleBox = styled.div.attrs({
  style: ({ color }) => ({
    backgroundColor: color
  })
})`
  width: 50px;
  padding: 50px;
`;

const LightState = ({ s, name }) => {
  if (s === undefined) return null;
  let { rgb, ...others } = s;
  return (
    <div>
      <SampleBox color={rgbToHex(rgb.r, rgb.g, rgb.b)} />
      <pre>{JSON.stringify(others)}</pre>
    </div>
  );
};
export default LightState;
