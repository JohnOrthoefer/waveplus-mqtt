package main

import (
   "fmt"
   "log"
   "strings"
   "encoding/json"
   "time"
   "github.com/eclipse/paho.mqtt.golang"
)

// Strings because the requirements for Homebridge
// This is the whole air quality that is supported by MQTTThing. 
// Wave Plus has Radon, Pariculate Matter (PM2.5), Volatile Organic Compounds (VOC), Carbon Dioxide (CO2), Temperature, and Humidity
// Overall AirQuality is based on wave numbers
//
// Implements Radon as Leak
// https://github.com/michaelahern/homebridge-airthings
type (
   AirQualityStruct struct {
      Timestamp               string `json:"timestamp, omitempty"`
      // (Unknown, excellent, good, fair, inferior, poor)
      OverallAirQuality       string `json:"airquality,omitempty"`

      // floating point number in units of parts per million
      CarbonDioxideLevel      float32 `json:"carbondioxidelevel,omitempty"`
      CarbonDioxideQuality    string  `json:"carbondioxidequality,omitempty"`

      // micrograms per cubic meter
      PM10Density             float32 `json:"pm10density,omitempty"`

      // micrograms per cubic meter
      PM2_5Density            float32 `json:"pm2_5density,omitempty"`

      // micrograms per cubic meter
      OzoneDensity            float32 `json:"ozonedensity,omitempty"`

      // micrograms per cubic meter
      NitrogenDioxideDensity  float32 `json:"nitrogendioxidedensity,omitempty"`

      // micrograms per cubic meter
      SulphurDioxideDensity   float32 `json:"sulphurdioxidedensity,omitempty"`

      // micrograms per cubic meter
      VOCDensity              float32 `json:"vocdensity,omitempty"`
      VOCQuality              string `json:"vocquality,omitempty"`

      // floating point number in units of parts per million
      CarbonMonoxideLevel     float32 `json:"carbonmonoxidelevel,omitempty"`
      CarbonMonoxideQuality   string  `json:"carbonmonoxidequality,omitempty"`

      // micrograms per cubic meter
      AirQualityPPM           float32 `json:"airqualityppm,omitempty"`

      StatusActive            bool   `json:"active,omitempty"`
      StatusFault             bool   `json:"fault,omitempty"` 
      StatusTampered          bool   `json:"tampered,omitempty"`
      StatusLowBattery        bool   `json:"lowbattery,omitempty"`

      // floating point number in degrees Celsius
      CurrentTemperature      float32 `json:"temperature,omitempty"`

      // floating point percentage representing the current relative humidity.
      CurrentRelativeHumidity float32 `json:"relativehumidit,omitempty"`

      // Not in Homekit/HomeBridge
      RadonShortTerm       float32  `json:"radonshort.omitempty`
      RadonShortQuality    string   `json:"radonshortquality.omitempty`
      RadonLongTerm        float32  `json:"radonlong.omitempty`
      RadonLongQuality     string   `json:"radonlongquality.omitempty`
   }
   wavePlusMQTT struct {
      m mqtt.Client
   }
)

func vocPPM2mgPm3(ppm float32) float32 {
   return ppm * 0.0409 * 100
}

func deg2C(degF string)string {
   c := (getFloat(degF)-32.0) * 5/9
   if c < 0 {
      c = 0
   }
   if c > 100 {
      c = 100
   }
   return fmt.Sprintf("%.1f", c)
}

func inHg2hPa(inHgStr string)string {
   inHg := getFloat(inHgStr) 
   hPa := int64(inHg * 33.86389)
   if hPa < 700 {
      hPa = 700
   }
   if hPa > 1100 {
      hPa = 1100 
   }
   
   return fmt.Sprintf("%d", hPa)
}

func newMQTT(broker string) *wavePlusMQTT {
   log.Printf("Broker: %s", broker)
   opts := mqtt.NewClientOptions().AddBroker(broker)
   mqttClient := wavePlusMQTT {
      m:  mqtt.NewClient(opts),
   }
   if token := mqttClient.m.Connect(); token.Wait() && token.Error() != nil {
      log.Printf("%s\n", token.Error())
      return nil
   }
   return &mqttClient
}

func batteryLow(s string) bool {
   if strings.ToLower(s) == "low" {
      return true
   }
   return false
}

func inch2mm(inch string) string {
   mmeters := getFloat(inch) * 25.4
   return fmt.Sprintf("%.1f", mmeters)
}

func mph2kph(mph string) string {
   kph := getFloat(mph) * 1.609
   return fmt.Sprintf("%.0f", kph)
}

func (m *wavePlusMQTT) publish(v *waveplus) {
   if m.m  == nil {
      return
   }
   AirQuality := AirQualityStruct{
      Timestamp: v.getTimestamp().Format(time.RFC822),
      OverallAirQuality: v.data.Quality().HKString(),
      CarbonDioxideLevel: v.data.co2Lvl,
      CarbonDioxideQuality: v.data.co2Quality().HKString(),
      VOCDensity: vocPPM2mgPm3(v.data.vocLvl),
      VOCQuality: v.data.vocQuality().HKString(),
      CurrentTemperature: v.data.temperature,
      CurrentRelativeHumidity: v.data.humidity,
      RadonShortTerm: v.data.radonShort,
      RadonShortQuality: v.data.radonShortQuality().HKString(),
      RadonLongTerm: v.data.radonLong,
      RadonLongQuality: v.data.radonLongQuality().HKString(),
   }

   jsonOut, _ := json.Marshal(AirQuality)
   //fmt.Printf("%s\n", jsonOut)
   //log.Printf("%d: Publishing to %s", v.sn, v.getMQTTTopic())
   token := m.m.Publish(fmt.Sprintf("%s/sample", v.getMQTTTopic()), 0, false, jsonOut)
   token.Wait()
}
