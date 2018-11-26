export class Light {
  static TypeHue = "hue";
  static TypeDMX = "DMX";
  static TypeGeneric = "generic";

  data: object;
  constructor(data: object) {
    this.data = data;
  }
  getType() {
    if ("hue_id" in this.data) return Light.TypeHue;
    if ("universe" in this.data) return Light.TypeDMX;
    return Light.TypeGeneric;
  }
}
