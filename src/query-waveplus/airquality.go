package main

func radon(v float32) AirQuality {
   v = convert2pCiL(v)

   if v < 2.7 {
      return Good
   }
   if v < 4.0 {
      return Fair
   }
   if v >= 4.0 {
      return Poor
   }

   return Unknown
}

func (w waveplusData)radonShortQuality() AirQuality {
   if !w.valid {
      return Unknown
   }
   return radon(w.radonShort)
}

func (w waveplusData)radonLongQuality() AirQuality {
   if !w.valid {
      return Unknown
   }
   return radon(w.radonLong)
}

func (w waveplusData)vocQuality() AirQuality {
   if !w.valid {
      return Unknown
   }

   if w.vocLvl < 250 {
      return Good
   }
   if w.vocLvl < 2000 {
      return Fair
   }
   if w.vocLvl >= 2000 {
      return Poor
   }
   return Unknown
}

func (w waveplusData)co2Quality() AirQuality {
   if !w.valid {
      return Unknown
   }

   if w.co2Lvl < 800 {
      return Good
   }
   if w.co2Lvl < 1000 {
      return Fair
   }
   if w.co2Lvl >= 1000 {
      return Poor
   }
   return Unknown
}

func (w waveplusData)humidityQuality() AirQuality {
   if !w.valid {
      return Unknown
   }

   v := w.humidity

   if v >=30 && v < 60 {
      return Good
   }
   if (v >= 60 && v < 70) || (v >= 25 && v < 30) {
      return Fair
   }
   if v >= 70 || v < 25 {
      return Poor
   }
   return Unknown
}

func (w waveplusData)temperatureQuality() AirQuality {
   if !w.valid {
      return Unknown
   }

   if w.temperature < 64 {
      return Good
   }
   if w.temperature >= 64 && w.temperature < 77 {
      return Fair
   }
   if w.temperature >= 77 {
      return Poor
   }
   return Unknown
}

func (w waveplusData)Quality() AirQuality {
   if !w.valid {
      return Unknown
   }

   

   return Good
}

func (a AirQuality)String() string {
   switch a {
   case Poor:
      return "Poor"
   case Fair:
      return "Fair"
   case Good:
      return "Good"
   }
   return "unknown"
}

