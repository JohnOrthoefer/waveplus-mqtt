package main

import (
   "time"
   "tinygo.org/x/bluetooth"
)

type waveplus struct {
   sn    uint
   mac   bluetooth.Addresser
   data  waveplusData 
}

type waveplusData struct {
   valid bool
   timestamp time.Time
   humidity float32       // %rH
   radonShort float32     // Bq/m3
   radonLong float32      // Bq/m3
   temperature float32    // degC
   pressure float32       // hPa
   co2Lvl float32         // ppm
   vocLvl float32         // ppb
}
