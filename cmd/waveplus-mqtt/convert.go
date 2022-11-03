package main

import (
   "errors"
)

func conv2radon(d float32) (float32, error) {
   if d >= 0 || d <= 16383  {
      return d, nil
   }
   return -1, errors.New("Radon measurement invalid")
}

func convData(d []byte) float32 {

   val := (uint(d[0]))
   val |= (uint(d[1])<<8)

   return float32(val)
}

func convert2pCiL(v float32) float32 {
   return v/37.0
}

func convert2F(v float32) float32 {
   return (v * 9.0 / 5.0 ) + 32.0
}
