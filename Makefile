all:  waveplus-mqtt install

waveplus-mqtt:
	$(MAKE) -C cmd/waveplus-mqtt

install:
	-mkdir bin
	cp cmd/waveplus-mqtt/waveplus_mqtt bin/

clean:
	$(MAKE) -C cmd/waveplus-mqtt clean
	rm bin/*
