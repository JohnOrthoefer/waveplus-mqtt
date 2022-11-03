package main
import (
//   "log"
)

// in pCi/L
func radon(v float32) AirQualityType {
   v = convert2pCiL(v)

   if v < 2.7 {
      return Excellent
   }
   if v < 4.0 {
      return Fair
   }
   if v >= 4.0 {
      return Poor
   }

   return Unknown
}

// in pCi/L
func (w waveplusData)radonShortQuality() AirQualityType {
   if !w.valid {
      return Unknown
   }
   return radon(w.radonShort)
}

// in pCi/L
func (w waveplusData)radonLongQuality() AirQualityType {
   if !w.valid {
      return Unknown
   }
   return radon(w.radonLong)
}

// in ppb
func (w waveplusData)vocQuality() AirQualityType {
   if !w.valid {
      return Unknown
   }

   if w.vocLvl < 250 {
      return Excellent
   }
   if w.vocLvl < 2000 {
      return Fair
   }
   if w.vocLvl >= 2000 {
      return Poor
   }
   return Unknown
}

// in ppm
func (w waveplusData)co2Quality() AirQualityType {
   if !w.valid {
      return Unknown
   }

   if w.co2Lvl < 800 {
      return Excellent
   }
   if w.co2Lvl < 1000 {
      return Fair
   }
   if w.co2Lvl >= 1000 {
      return Poor
   }
   return Unknown
}

// in %rH
func (w waveplusData)humidityQuality() AirQualityType {
   if !w.valid {
      return Unknown
   }

   v := w.humidity

   if v >=30 && v < 60 {
      return Excellent
   }
   if (v >= 60 && v < 70) || (v >= 25 && v < 30) {
      return Fair
   }
   if v >= 70 || v < 25 {
      return Poor
   }
   return Unknown
}

// in F
func (w waveplusData)temperatureQuality() AirQualityType {
   if !w.valid {
      return Unknown
   }

   if w.temperature < 64 {
      return Excellent
   }
   if w.temperature >= 64 && w.temperature < 77 {
      return Fair
   }
   if w.temperature >= 77 {
      return Poor
   }
   return Unknown
}

func (w waveplusData)Quality() AirQualityType {
   if !w.valid {
      return Unknown
   }

   rtn := w.radonShortQuality().AirQualityComp(w.vocQuality())
   rtn  = rtn.AirQualityComp(w.co2Quality())
   //rtn  = rtn.AirQualityComp(w.humidityQuality())
   //rtn  = rtn.AirQualityComp(w.temperatureQuality())
   
   return rtn
}

func (a AirQualityType)AirQualityComp(b AirQualityType) AirQualityType {
   if a == Unknown || b == Unknown { return Unknown }
   if a == b { return a }
   switch a {
   case Poor: 
      return Poor
   case Inferior:
      if b == Poor { return Poor }
      return Inferior
   case Fair:
      if b == Poor { return Poor }
      if b == Inferior { return Inferior }
      return Fair
   case Good:
      if b == Poor { return Poor }
      if b == Inferior { return Inferior }
      if b == Fair { return Fair }
      return Good
   case Excellent:
      if b == Poor { return Poor }
      if b == Inferior { return Inferior }
      if b == Fair { return Fair }
      if b == Good { return Good }
      return Excellent
   }
   return Unknown
}

func (a AirQualityType)String() string {
   switch a {
   case Poor:
      return "Poor"
   case Inferior:
      return "Inferior"
   case Fair:
      return "Fair"
   case Good:
      return "Good"
   case Excellent:
      return "Excellent"
   }
   return "unknown"
}

func (a AirQualityType)HKString() string {
   switch a {
   case Poor:
      return "POOR"
   case Inferior:
      return "INFERIOR"
   case Fair:
      return "FAIR"
   case Good:
      return "GOOD"
   case Excellent:
      return "EXCELLENT"
   }
   return "UNKNOWN"
}
