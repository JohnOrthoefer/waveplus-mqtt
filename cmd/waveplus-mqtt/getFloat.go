package main

import (
   "strconv"
)

// turn a string into Float64
func getFloat(s string) float64 {
   v, err := strconv.ParseFloat(s, 64)
   if err != nil {
      nan, _ := strconv.ParseFloat("NaN", 64)
      v = nan
   }
   return v
}
