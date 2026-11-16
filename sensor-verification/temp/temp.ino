#include <OneWire.h>

// Data wire is connected to GPIO 4
#define ONE_WIRE_BUS 4

OneWire oneWire(ONE_WIRE_BUS);

uint8_t addr[8];

void setup() {
  Serial.begin(9600);

  delay(100);

  oneWire.reset_search();
  bool found = oneWire.search(addr);
  if (!found) {
    Serial.println("Device not found");
    while (true) ;;
  }
  Serial.print("model id: "); Serial.println(addr[0], HEX); // 0x28 == DS18B20
  // Serial.print("id: "); Serial.println(addr);

  // read powersupply
  oneWire.reset();
  oneWire.select(addr);
  oneWire.write(0xB4);
  uint8_t mode = oneWire.read_bit();
  oneWire.reset();
  Serial.print("power mode: "); Serial.println(mode);


  // scratchpad: 8 bytes. 2-4 writable
  // [temp LSB, temp MSB], [TH/user1, TL/user2], conf, _reserved_, _reserved_, _reserved_, CRC

  // configure
  // 12-bit resolution mode, no alarms
  oneWire.reset();
  oneWire.select(addr);
  oneWire.write(0x4E);
  oneWire.write(0xAA); // High alarm, signed 8bit
  oneWire.write(0x55); // low alarm, signed 8bit
  oneWire.write(0x03 << 5 | 0b00011111); // config. Only resolution (bits 6 and 5) are settable. bits 4-0 are always 1
  oneWire.reset();

  // read back configuration and verify CRC
  uint8_t scratch[9];
  oneWire.reset();
  oneWire.select(addr);
  oneWire.write(0xBE);
  oneWire.read_bytes(scratch, 9);
  oneWire.reset();

  Serial.print("scratch: ");
  for (size_t i = 0; i<9; i++) {
    Serial.print(scratch[i], HEX);
    Serial.print(" ");
  }
  Serial.println();

  // verify CRC
  uint8_t crc = oneWire.crc8(scratch, 8);
  Serial.print("crc: "); Serial.println(crc, HEX);
  if (crc != scratch[8]) {
    Serial.println("CRC mismatch");
    while (true) ;;
  }
}

void loop() {
  oneWire.reset();
  oneWire.select(addr);
  oneWire.write(0x44); // start conversion

  // in 12bit mode, conversion time takes 750ms
  delay(750);
  // wait for complete
  while (true) {
    if (oneWire.read_bit()) {
      break;
    }
  }

  uint8_t scratch[9];
  oneWire.reset();
  oneWire.select(addr);
  oneWire.write(0xBE);
  oneWire.read_bytes(scratch, 9);
  oneWire.reset();

  Serial.print("scratch: ");
  for (size_t i = 0; i<9; i++) {
    Serial.print(scratch[i], HEX);
    Serial.print(" ");
  }
  Serial.println();

  // verify CRC
  uint8_t crc = oneWire.crc8(scratch, 8);
  Serial.print("crc: "); Serial.println(crc, HEX);
  if (crc != scratch[8]) {
    Serial.println("CRC mismatch");
    while (true) ;;
  }

  // TODO: is this true??
  if (scratch[0] == 0x50 && scratch[1] == 0x05 && scratch[6] == 0x0C) {
    Serial.println("power on reset");
    delay(1000);
    return;
  }

  if (scratch[0] == 0xFF && scratch[1] == 0x07) {
    Serial.println("insufficent power"); // shouldn't happen outside of parasite mode
    while (true) ;;
  }

  // int32_t fpTemperature = 0;
  // // looking thru the spec sheets of all supported devices, bit 15 is always the signing bit
  // int32_t neg = 0x0;
  // if (scratch[0] & 0x80)
  //     neg = 0xFFF80000;
  // // We must mask out bit 1 (reserved) and 0 (fault) on TEMP_LSB
  // fpTemperature = (((int32_t)scratch[0]) << 11)
  //                 | (((int32_t)scratch[1] & 0xFC) << 3)
  //                 | neg;

  Serial.print("sign: "); Serial.println(scratch[1] >> 3, BIN);
  int16_t tempC = scratch[0] | ((scratch[1] & 0x03) << 8);
  if (scratch[1] & (1 << 7)) {
    tempC = -1*tempC;
  }

  Serial.print("temp raw: "); Serial.println(tempC, DEC);
  Serial.print("temp C: "); Serial.println((float)tempC/16.0);

  delay(100);

  //  sensors.requestTemperatures(); 
   
  //  delay(750); 
   
  //  float tempC = sensors.getTempCByIndex(0);
  // //  Serial.print("Temperature: ");
  //  Serial.println(tempC);
  // //  Serial.println("°C");
  //  delay(100);
}

