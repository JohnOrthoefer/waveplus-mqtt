package main

import (
   "os"
   "log"
   "time"
)

func main() {
   var mon []*waveplus

   for _, v := range os.Args[1:] {
      mon = append(mon, newMonitor(v))
   }

   c := ReadYAML()

   log.Printf("yaml: %q", c)
   for _, v := range c.Monitors() {
      log.Printf("Serial: %s", v.SerialNumber())
      mon = append(mon, newMonitor(v.SerialNumber()))
   }

   if len(mon) < 1 {
      log.Fatal("Nothing to monitor")
   }

   waitTime := time.Second * 30

   for {
      thisRun := time.Now()
      for _, v := range mon {
         if !v.ready() {
            v.getMonitorMAC(time.Second * 30)
         }
         v.getMonitorValues()
         v.printMonitorValues()
      }
      time.Sleep(time.Until(thisRun.Add(waitTime)))
   }
}
