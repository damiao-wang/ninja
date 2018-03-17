package main

import (
	"ninja/base/misc/cmdt"
	"net/http"

	"github.com/getsentry/raven-go"
)

func main() {
	cmdt.SetName("gunbuster")
	 s := http.Server
	 s.Shutdown()
}
