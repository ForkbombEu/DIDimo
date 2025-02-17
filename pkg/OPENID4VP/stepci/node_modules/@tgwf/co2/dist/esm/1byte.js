const CO2_PER_KWH_IN_DC_GREY = 519;
const CO2_PER_KWH_NETWORK_GREY = 475;
const CO2_PER_KWH_IN_DC_GREEN = 0;
const KWH_PER_BYTE_IN_DC = 72e-12;
const FIXED_NETWORK_WIRED = 429e-12;
const FIXED_NETWORK_WIFI = 152e-12;
const FOUR_G_MOBILE = 884e-12;
const KWH_PER_BYTE_FOR_NETWORK = (FIXED_NETWORK_WIRED + FIXED_NETWORK_WIFI + FOUR_G_MOBILE) / 3;
const KWH_PER_BYTE_FOR_DEVICES = 13e-11;
class OneByte {
  constructor(options) {
    this.options = options;
    this.KWH_PER_BYTE_FOR_NETWORK = KWH_PER_BYTE_FOR_NETWORK;
  }
  perByte(bytes, green) {
    if (bytes < 1) {
      return 0;
    }
    if (green) {
      const Co2ForDC = bytes * KWH_PER_BYTE_IN_DC * CO2_PER_KWH_IN_DC_GREEN;
      const Co2forNetwork = bytes * KWH_PER_BYTE_FOR_NETWORK * CO2_PER_KWH_NETWORK_GREY;
      return Co2ForDC + Co2forNetwork;
    }
    const KwHPerByte = KWH_PER_BYTE_IN_DC + KWH_PER_BYTE_FOR_NETWORK;
    return bytes * KwHPerByte * CO2_PER_KWH_IN_DC_GREY;
  }
}
var byte_default = OneByte;
export {
  OneByte,
  byte_default as default
};
