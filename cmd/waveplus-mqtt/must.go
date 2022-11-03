package main

import (
   "log"
)

func must(action string, err error) {
	if err != nil {
		log.Fatal("failed to " + action + ": " + err.Error())
	}
}
