package utils

import (
	"encoding/json"
	"log"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func ParseBody(r []byte, x interface{}) error {
	if err := json.Unmarshal(r, x); err != nil {
		return err
	}
	return nil
}
