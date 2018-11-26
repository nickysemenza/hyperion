export const LightTypeHue = "hue";
export const LightTypeDMX = "DMX";
export const LightTypeGeneric = "generic";

export const getLightType = (light: object) => {
  if ("hue_id" in light) return LightTypeHue;
  if ("universe" in light) return LightTypeDMX;
  return LightTypeGeneric;
};

const componentToHex = (c: number) => {
  var hex = c.toString(16);
  return hex.length === 1 ? "0" + hex : hex;
};

export const rgbToHex = (
  r: number,
  g: number,
  b: number,
  withoutPound = false
) =>
  (withoutPound ? "" : "#") +
  componentToHex(r) +
  componentToHex(g) +
  componentToHex(b);

export const isRGBDark = (r: number, g: number, b: number) =>
  //calculate perceptive luminance
  r * 0.299 + g * 0.587 + b * 0.114 / 255 > 0.5;
