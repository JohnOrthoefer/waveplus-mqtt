package main

import (
   "os"
   "log"
   "time"
)

func main() {
   var mon []*waveplus

   if len(os.Args) < 2 {
      log.Fatal("no params")
   }
   
   for _, v := range os.Args[1:] {
      mon = append(mon, newMonitor(v))
   }

   for {
      for _, v := range mon {
         if !v.ready() {
            v.getMonitorMAC(time.Second * 30)
         }
         v.getMonitorValues()
         v.printMonitorValues()
      }
      time.Sleep(time.Second *10)
   }
}
