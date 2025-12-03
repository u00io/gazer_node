package main

import (
	"github.com/u00io/gazer_node/forms/mainform"
	"github.com/u00io/gazer_node/localstorage"
	"github.com/u00io/gomisc/logger"
)

func main() {
	localstorage.Init("gazer_node")
	logger.Init(localstorage.Path() + "/logs")
	mainform.Run()
}
