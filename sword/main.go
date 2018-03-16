package main

import (
	"ninja/base/misc/cmdt"
	"net/http"

	"github.com/getsentry/raven-go"
)

func init() {
	raven.SetDSN("https://1c115aaf0b2048f485936409b03ce0f7:c6480693facf40d6b89c28711fb363b9@sentry.io/304312")
}

func main() {
	cmdt.SetName("gunbuster")
	 s := http.Server
	 s.Shutdown()
}
