package main

import (
	"github.com/KubaiDoLove/conductor/internal/app/apiserver"
)

var version = "unknown"

func main() {
	apiserver.Start(version)
}
