void setup() {
  Serial.begin(9600);
}

void loop() {
  uint16_t v = analogRead(A0);
  Serial.println(v);
  delay(100);
}

