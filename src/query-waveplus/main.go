package main

import (
   "os"
   "time"
)

func main() {

   mon := getMonitors(os.Args[1:])
   
   for _, v := range mon {
      v.getMonitorMAC()
   }

   for {
      for _, v := range mon {
         v.getMonitorValues()
         v.printMonitorValues()
      }
      time.Sleep(time.Second *10)
   }
}
