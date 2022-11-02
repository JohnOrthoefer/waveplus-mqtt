package main

import (
   "fmt"
   "strings"
   "encoding/json"
   "net/url"
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
      // (Unknown, excellent, good, fair, inferior, poor)
      OverallAirQuality       string `json:"airquality,omitempty"`

      // floating point number in units of parts per million
      CarbonDioxideLevel      float64 `json:"carbondioxidelevel,omitempty"`

      // micrograms per cubic meter
      PM10Density             float64 `json:"pm10density,omitempty"`

      // micrograms per cubic meter
      PM2_5Density            float64 `json:"pm2_5density,omitempty"`

      // micrograms per cubic meter
      OzoneDensity            float64 `json:"ozonedensity,omitempty"`

      // micrograms per cubic meter
      NitrogenDioxideDensity  float64 `json:"nitrogendioxidedensity,omitempty"`

      // micrograms per cubic meter
      SulphurDioxideDensity   float64 `json:"sulphurdioxidedensity,omitempty"`

      // micrograms per cubic meter
      VOCDensity              float64 `json:"vocdensity,omitempty"`

      // floating point number in units of parts per million
      CarbonMonoxideLevel     float64 `json:"carbonmonoxidelevel,omitempty"`

      // micrograms per cubic meter
      AirQualityPPM           float64 `json:"airqualityppm,omitempty"`
      StatusActive            bool   `json:"active,omitempty"`
      StatusFault             bool   `json:"fault,omitempty"` 
      StatusTampered          bool   `json:"tampered,omitempty"`
      StatusLowBattery        bool   `json:"lowbattery,omitempty"`

      // floating point number in degrees Celsius
      CurrentTemperature      string `json:"temperature,omitempty"`

      // floating point percentage representing the current relative humidity.
      CurrentRelativeHumidity float64 `json:"relativehumidit,omitempty"`
   }
)

var (
   AirQuality AirQualityType
   mqttClient mqtt.Client
   mqttTopic  string
)


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

func setTopic(t string) {
   mqttTopic = t
}

func getTopic(v url.Values) string {
   rtn := mqttTopic
/*
   for _, val := range []string{"ID", "softwaretype", "id", "mt", "sensor"} {
      rtn = strings.Replace(rtn, fmt.Sprintf("%%%s%%", val), v.Get(val), -1)
   }
*/
   //fmt.Printf("Topic: %s => %s\n", mqttTopic, rtn)
   return rtn
}

func mqttSetup(broker string, topic string) {
   opts := mqtt.NewClientOptions().AddBroker(broker)
   mqttClient = mqtt.NewClient(opts)
   if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
      mqttClient = nil
      fmt.Printf("%s\n", token.Error())
      return
   }
   setTopic(topic)
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

func publish(v url.Values) {
   if mqttClient == nil {
      return
   }
   AirQuality := AirQualityStruct{
/*
      Temperature:   deg2C(v.Get("tempf")),
      Humidity:      v.Get("humidity"),
      AirPressure:   inHg2hPa(v.Get("baromin")),
      Rain1h:        inch2mm(v.Get("rainin")),
      Rain24h:       inch2mm(v.Get("dailyrainin")),
      WindDirection: v.Get("winddir"),
      WindSpeed:     mph2kph(v.Get("windspeedmph")),
      Battery:       batteryLow(v.Get("sensorbattery")),
*/
   }

   jsonOut, _ := json.Marshal(AirQuality)
   //fmt.Printf("%s\n", jsonOut)
   mqttClient.Publish(getTopic(v), 0, false, jsonOut)
}
