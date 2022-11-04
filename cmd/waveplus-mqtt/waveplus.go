package main

import (
   "fmt"
   "log"
   "time"
   "strings"
   "strconv"
   "tinygo.org/x/bluetooth"
)

// 
// 
func newMonitor(m string) *waveplus {
   var rtn waveplus

   if len(m) == 0 {
      log.Fatal("Needs a Serial Number")
   }

   if s, err := strconv.ParseUint(m, 10, 32); err == nil {
      rtn.sn = uint(s)
   } else {
      return nil
   }
   
   rtn.location = m

   return &rtn
}

func (m *waveplus) setLocation(l string) {
   if l == "" {
      return
   }
   m.location = l
}

func (m *waveplus) getLocation() string {
   return m.location
}

func getLocations(m []*waveplus) string {
   var rtn []string
   for _, v := range m {
      rtn = append(rtn, v.getLocation())
   }
   return strings.Join(rtn, ", ")
}

func (m *waveplus) getSerialNumber() uint {
   return m.sn
}

func startScan() {
   var adapter = bluetooth.DefaultAdapter

   must("enable BLE stack", adapter.Enable())
	go adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {})
}

func (w *waveplus) getMonitorMAC(wait time.Duration) {
   var adapter = bluetooth.DefaultAdapter

   must("enable BLE stack", adapter.Enable())

   //log.Printf("%d: Searching for monitor, %s", w.sn, w.location)

   timeout := time.Now().Add(wait)
	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
      if v, ok := (result.ManufacturerData())[0x0334]; ok {
         // log.Printf("Found ManufacturerData, (%d == %d)\n", getSerialNumber(v[0:4]), w.sn)
         if getSerialNumber(v[0:4]) ==  w.sn {
            w.mac = result.Address
			   adapter.StopScan()
         }
         if time.Until(timeout) < 0 {
			   adapter.StopScan()
         }
		}
	})
   if err != nil {
      log.Printf("Scan Fail %s", err)
   }
}

func (w *waveplus) ready() bool {
   if w.mac == nil {
      return false
   }
   return true
}

func (w *waveplus) getMonitorValues() bool {
   if !w.ready() {
      log.Printf("%d: Not ready", w.sn)
      return false
   }

   //log.Printf("%d: quering\n", w.sn)
   w.data.valid = false

   var adapter = bluetooth.DefaultAdapter

   must("enable BLE stack", adapter.Enable())

   //log.Printf("%d: MAC %s", w.sn, w.mac.String())

	device, err := adapter.Connect(w.mac, bluetooth.ConnectionParams{})
   if err != nil {
      log.Printf("%d: Failed to Connect, %s", w.sn, err.Error())
      return false
   }

	srvcs, err := device.DiscoverServices(nil)
   if err != nil {
      log.Printf("Failed discover service %s\n", err)
      return false
   }

   for _, srvc := range srvcs {
      chars, err := srvc.DiscoverCharacteristics(nil)
      must("Discover Characteristics", err)
      for _, char := range chars {
         if char.UUID().String() == "b42e2a68-ade7-11e4-89d3-123b93f75cba" {
            buf := make([]byte, 255)
            n, err := char.Read(buf)
            if err != nil || n < 20 {
               log.Printf("Invalid Characterstics Read")
               continue
            }
            w.data = updateData(buf)
            if w.data.valid {
               w.samples += 1
            }
         }
      }
   }

   device.Disconnect()
   return w.data.valid
}

func (w *waveplus) printMonitorSummery() {
   if !w.data.valid {
      return
   }
   var rtn = []string {
      w.data.Quality().String(),
      fmt.Sprintf("%.1f", convert2pCiL(w.data.radonShort))+"pCi/L",
      fmt.Sprintf("%.1f", convert2pCiL(w.data.radonLong))+"pCi/L",
      fmt.Sprintf("%.0f", w.data.vocLvl)+"ppb",
      fmt.Sprintf("%.0f", w.data.co2Lvl)+"ppm",
      fmt.Sprintf("%.1f", w.data.humidity)+"%rH",
      fmt.Sprintf("%.1f", convert2F(w.data.temperature))+"F",
      fmt.Sprintf("%.0f", w.data.pressure)+"hPa",
   }
   log.Printf("%s/%d: %s\n", w.getLocation(), w.sn, strings.Join(rtn, ","))
}

func (w *waveplus) printMonitorValues() {
   if !w.data.valid {
      return
   }
   log.Printf("%d: Radon- Short %3.1f pCi/L(%s), Long %3.1f pCi/L(%s)\n",
      w.sn, 
      convert2pCiL(w.data.radonShort), w.data.radonShortQuality().String(),
      convert2pCiL(w.data.radonLong), w.data.radonLongQuality().String())
   log.Printf("%d: VOC- %4.0f ppb(%s) CO2- %4.0f ppm(%s)\n", 
      w.sn, 
      w.data.vocLvl, w.data.vocQuality().String(),
      w.data.co2Lvl, w.data.co2Quality().String())

   log.Printf("%d: %3.1f %%rH, %5.1f F, %4.0f hPa\n", 
      w.sn, 
      w.data.humidity, 
      convert2F(w.data.temperature),
      w.data.pressure)
   log.Printf("%d: %s, %d samples\n", w.sn, w.data.Quality().String(), w.samples)
}

func (w *waveplus) setMQTTTopic(s string) {
   if s == "" {
      s = fmt.Sprintf("tele/%d", w.sn)
   }
   w.mqttTopic = s
}

func (w *waveplus) getMQTTTopic() string {
   if w.mqttTopic == "" {
      w.setMQTTTopic("")
   }
   return w.mqttTopic
}

func updateData(d []byte) waveplusData {
   var rtn waveplusData 

   // Assume everything is valid
   rtn.valid = true
   rtn.timestamp = time.Now()

   // Byte 0: version 1 data only
   if d[0] != 1 {
      rtn.valid = false
      return rtn
   }

   // Byte 1: Humidity 
   rtn.humidity = float32((uint(d[1]))/2.0)

   // Byte 4-5: Radon Short Term
   rs, err := conv2radon(convData(d[4:6]))
   if err != nil {
      rtn.valid = false
   }
   rtn.radonShort = rs

   // Radon Long Term
   rl, err := conv2radon(convData(d[6:8]))
   if err != nil {
      rtn.valid = false
   }
   rtn.radonLong = rl

   // Temperature
   rtn.temperature = convData(d[8:10]) /100.0

   // Atomsphere Pressure
   rtn.pressure = convData(d[10:12]) / 50.0

   // CO2 Level
   rtn.co2Lvl = convData(d[12:14]) 

   // Vol Organic Chem  Level
   rtn.vocLvl = convData(d[14:16]) 

   return rtn
}
