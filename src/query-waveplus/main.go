package main

import (
   "os"
   "log"
   "time"
)

func main() {
   var mon []*waveplus

   log.Printf("# %s - %s @%s\n", repoName, sha1ver, buildTime)

   for _, v := range os.Args[1:] {
      mon = append(mon, newMonitor(v, v))
   }

   c := ReadYAML()

   if !c.TimeStampsEnabled() {
      log.SetFlags(0)
   }

   for _, v := range c.Monitors() {
      mon = append(mon, newMonitor(v.SerialNumber(), v.GetLocation()))
   }

   if len(mon) < 1 {
      log.Fatal("Nothing to monitor")
   }

   for {
      thisRun := time.Now()
      for _, v := range mon {
         if !v.ready() {
            v.getMonitorMAC(c.GetTimeout())
         }
         v.getMonitorValues()
         v.printMonitorValues()
      }
      time.Sleep(time.Until(thisRun.Add(c.GetFrequency())))
   }
}
