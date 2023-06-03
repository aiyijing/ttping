package tt

import (
	"fmt"
	"os"
)

type Pinger interface {
	Ping()
}

func NewPinger(cfg *Config) Pinger {
	var pinger Pinger
	switch cfg.Protocol {
	case TCP:
		pinger = NewTCPinger(cfg)
	case HTTP:
		pinger = NewHTTPPinger(cfg)
	case ICMP:
		pinger = NewICMPPinger(cfg)
	default:
		fmt.Println("invalid protocol type")
		os.Exit(1)
	}
	return pinger
}
