package service

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

func LogError(err error, message string) bool {
	if err != nil {
		log.Err(err).Msg(fmt.Sprint(err, message))
		return true
	}
	return false
}
