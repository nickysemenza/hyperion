import React from 'react';
import styled from 'styled-components';
import { rgbToHex } from '../utils';

const SampleBox = styled.div.attrs({
  style: ({ color }) => ({
    backgroundColor: color
  })
})`
  width: 20px;
  padding: 20px;
`;

export const ColorBox = ({ state }) => {
  if (state === undefined) return null;
  let { rgb } = state;
  return <SampleBox color={rgbToHex(rgb.r, rgb.g, rgb.b)} />;
};

export const LightState = ({ s, name }) => {
  if (s === undefined) return null;
  let { rgb, ...others } = s;
  return (
    <div>
      <SampleBox color={rgbToHex(rgb.r, rgb.g, rgb.b)} />
      <code>{JSON.stringify(others)}</code>
    </div>
  );
};
