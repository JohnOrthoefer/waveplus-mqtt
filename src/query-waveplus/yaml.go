package main

import (
   "log"
   "fmt"
   "os"
   "time"
   "gopkg.in/yaml.v3"
)

type Configuration struct {
   Freq     uint              `yaml:"frequency"`
   Timeout  uint              `yaml:"timeout"`
   TSEnabled bool             `yaml:"timestamps"`
   MqttURL   string      `yaml:"mqtturl"`
   Monitor  []MonitorRecord   `yaml:"monitor"`
}

type MonitorRecord struct {
   Location string            `yaml:"location"`
   Serial   uint              `yaml:"serialnumber"`
   Mac      string            `yaml:"macaddr"`
   Topic    string            `yaml:"topic"`
}

func ReadYAML() *Configuration {
   conf := Configuration {
      Freq: 60,
      Timeout: 60, 
      TSEnabled: true,
      MqttURL: "",
      Monitor: nil,
   }

   yamlFile, err := os.ReadFile("waveplus.yaml")
   if err != nil {
      log.Printf("YAML Configuration, %s", err)
      return nil
   }

   err = yaml.Unmarshal(yamlFile, &conf)
   if err != nil {
      log.Printf("YAML Parse Error, %s", err)
      return nil
   }

   return &conf
}

func (c *Configuration) Monitors() []MonitorRecord {
   return c.Monitor
}

func (c *Configuration) TimeStampsEnabled() bool {
   return c.TSEnabled
}

func (m MonitorRecord) SerialNumber() string {
   return fmt.Sprintf("%d", m.Serial)
}

func (m MonitorRecord) GetMqttTopic() string {
   if m.Topic == "" {
      return fmt.Sprintf("tele/%s", m.SerialNumber())
   }
   return m.Topic
}

func (m MonitorRecord) GetLocation() string {
   if m.Location == "" {
      return m.SerialNumber()
   }
   return m.Location
}

func (c *Configuration) GetFrequency() time.Duration {
   return time.Second * time.Duration(c.Freq)
}

func (c *Configuration) GetTimeout() time.Duration {
   return time.Second * time.Duration(c.Timeout)
}

func (c *Configuration) GetBroker() string {
   rtn := c.MqttURL
   if rtn == "" {
      return "tcp://localhost:1883"
   }
   return rtn
}

