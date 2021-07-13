package testbed

import (
	"github.com/rs/zerolog/log"
)

func showLog() {
	notify := false
	log.Info().Bool("resolution", notify).Msg("Should notify user") // zero log
	log.Print("hello")                                              // default go style log
}
