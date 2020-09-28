package engine

import (
	"github.com/Misora000/mangadl/engine/router"
	"github.com/Misora000/mangadl/logging"
	"github.com/Misora000/mangadl/parser"
	"github.com/Misora000/mangadl/types"
)

//Initialize initializes engine.
func Initialize() {
	parser.Initialize()
}

// Finalize finalizes engine.
func Finalize() {
}

// GetParser get parser from router.
func GetParser(URL string) types.Parser {
	if p := router.ParseURL(URL); p != nil {
		return p
	}
	logging.Fatal("no available parser")
	return nil
}

// PrintSupportedSites prints all supported sites.
func PrintSupportedSites() {
	logging.Log("Supported sites:")
	for _, v := range router.ListRoutingBook() {
		logging.Log("\t%v", v)
	}
}
