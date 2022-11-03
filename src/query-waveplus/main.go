package main

import (
   "os"
   "log"
   "time"
)

func main() {
   var monLst []*waveplus

   log.Printf("# %s - %s @%s\n", repoName, sha1ver, buildTime)

   for _, v := range os.Args[1:] {
      mon := newMonitor(v)
      monLst = append(monLst, mon)
   }

   c := ReadYAML()

   if !c.TimeStampsEnabled() {
      log.SetFlags(0)
   }

   broker := newMQTT(c.GetBroker())
   if broker == nil {
      log.Fatal("Error Setting up MQTT Broker")
   }

   for _, v := range c.Monitors() {
      mon := newMonitor(v.SerialNumber())
      mon.setLocation(v.GetLocation())
      mon.setMQTTTopic(v.GetMqttTopic())
      monLst = append(monLst, mon)
   }

   if len(monLst) < 1 {
      log.Fatal("Nothing to monitor")
   }

   for {
      thisRun := time.Now()
      for _, v := range monLst {
         if !v.ready() {
            v.getMonitorMAC(c.GetTimeout())
         }
         if v.getMonitorValues() {
            broker.publish(v)
         }
         v.printMonitorValues()
      }
      time.Sleep(time.Until(thisRun.Add(c.GetFrequency())))
   }
}
