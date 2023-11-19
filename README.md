## Airthings Waveplus Airquality bridge to MQTT

### Intro
This provides a way to get Air Quality data from a Waveplus to Homekit. 

### Building 
``` 
make all
```

### Install
```
sudo useradd --system --shell /sbin/nologin waveplus
sudo cp bin/waveplus_mqtt /usr/local/bin/
sudo cp etc/waveplus.yaml /usr/local/etc/
sudo chown waveplus:waveplus /usr/local/etc/waveplus.yaml
sudo cp etc/waveplus.server /etc/systemd/system/
sudo nano /usr/local/etc/waveplus.yaml
# add your 
# Serial Numbers required
$ optional locations, and the topic you want to publish to
sudo systemctl enable --now waveplus
```

