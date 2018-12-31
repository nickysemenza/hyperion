import React from 'react';
import styled from '../util/styled-components';

import { rgbToHex } from '../utils';

const SampleBox = styled('div')<{ color: string }>`
  background-color: ${props => props.color};
  width: 20px;
  padding: 20px;
`;

type ColorBoxProps = {
  state: any;
};
export const ColorBox: React.SFC<ColorBoxProps> = ({ state }) => {
  if (state === undefined) return null;
  let { rgb } = state;
  return <SampleBox color={rgbToHex(rgb.r, rgb.g, rgb.b)} />;
};

type LightStateProps = {
  s: any;
  name: string;
};
export const LightState: React.SFC<LightStateProps> = ({ s, name }) => {
  if (s === undefined) return null;
  let { rgb, ...others } = s;
  return (
    <div>
      <SampleBox color={rgbToHex(rgb.r, rgb.g, rgb.b)} />
      <code>{JSON.stringify(others)}</code>
    </div>
  );
};
