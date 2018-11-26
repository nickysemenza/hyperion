import { rgbToHex } from './utils';
import { Light } from './types';
test('correctly identifies light type', () => {
  let tt = [
    {
      light: {
        hue_id: 1,
        name: 'hue1'
      },
      expected: Light.TypeHue
    },
    {
      light: {
        name: 'par2',
        start_address: 8,
        universe: 1,
        profile: 'china-par-1'
      },
      expected: Light.TypeDMX
    },
    {
      light: {
        name: 'generic1'
      },
      expected: Light.TypeGeneric
    }
  ];
  tt.forEach(tc => {
    let l = new Light(tc['light']);
    expect(l.getType()).toBe(tc['expected']);
  });
});

test('converting rgb to hex', () => {
  expect(rgbToHex(255, 0, 0)).toBe('#ff0000');
  expect(rgbToHex(255, 0, 0, true)).toBe('ff0000');
  expect(rgbToHex(0, 0, 0)).toBe('#000000');
});
