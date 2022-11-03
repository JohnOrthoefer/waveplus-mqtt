package main

import (
   "strconv"
)

func atoui(v string) uint {
   if s, err := strconv.ParseUint(v, 10, 32); err == nil {
      return uint(s)
   }

   return 0
}
