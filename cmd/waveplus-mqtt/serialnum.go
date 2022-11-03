package main

func getSerialNumber(value []byte) uint {
	if len(value) > 3 {

		var sn1 uint

		sn1 = (uint)(value[0])
		sn1 |= ((uint)(value[1])) << 8
		sn1 |= ((uint)(value[2])) << 16
		sn1 |= ((uint)(value[3])) << 24
		return sn1
	}
	return 0
}
