all:  waveplus-mqtt install

waveplus-mqtt:
	$(MAKE) -C cmd/waveplus-mqtt

install:
	-mkdir bin
	cp cmd/waveplus-mqtt/waveplus_mqtt bin/

update:
	go get -u github.com/eclipse/paho.mqtt.golang 
	go get -u gopkg.in/yaml.v3 
	go get -u tinygo.org/x/bluetooth
	go mod tidy

clean:
	$(MAKE) -C cmd/waveplus-mqtt clean
	rm bin/*
