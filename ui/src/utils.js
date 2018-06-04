export const LightTypeHue = 'hue';
export const LightTypeDMX = 'DMX';
export const LightTypeGeneric = 'generic';

export const getLightType = light => {
  if ('hue_id' in light) return LightTypeHue;
  if ('universe' in light) return LightTypeDMX;
  return LightTypeGeneric;
};

const componentToHex = c => {
  var hex = c.toString(16);
  return hex.length === 1 ? '0' + hex : hex;
};

export const rgbToHex = (r, g, b, withoutPound = false) =>
  (withoutPound ? '' : '#') +
  componentToHex(r) +
  componentToHex(g) +
  componentToHex(b);
