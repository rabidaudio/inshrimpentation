#include <OneWire.h>
#include <DallasTemperature.h>

// Data wire is connected to GPIO 4
#define ONE_WIRE_BUS 4

OneWire oneWire(ONE_WIRE_BUS);
DallasTemperature sensors(&oneWire);


void setup() {
  Serial.begin(9600);
  sensors.begin();
}

void loop() {
   sensors.requestTemperatures(); 
   
   delay(750); 
   
   float tempC = sensors.getTempCByIndex(0);
  //  Serial.print("Temperature: ");
   Serial.println(tempC);
  //  Serial.println("°C");
   delay(100);
}

