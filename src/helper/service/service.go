package service

import (
	"fmt"
	"log"
)

func HandleError(err error, message string) bool {
	if err != nil {
		log.Println(fmt.Sprint(err, message))
		return true
	}
	return false
}
