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
