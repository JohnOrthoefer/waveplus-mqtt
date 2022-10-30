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

