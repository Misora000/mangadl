package router

import (
	"net/url"
	"sort"

	"github.com/Misora000/mangadl/logging"
	"github.com/Misora000/mangadl/types"
)

// ParserGenerator is the generator function of parser.
type ParserGenerator func() types.Parser

var routingBook = map[string]func() types.Parser{}

// RegisterParser register parser.
func RegisterParser(route string, f ParserGenerator) {
	if _, exist := routingBook[route]; exist {
		logging.Error("routing %v is already existed", route)
	}
	routingBook[route] = f
}

// ParseURL gets parser for the given url.
func ParseURL(URL string) types.Parser {
	u, err := url.Parse(URL)
	if err != nil {
		logging.Fatal(err.Error())
	}

	if generator, exist := routingBook[u.Host]; exist {
		return generator()
	}
	return nil
}

// ListRoutingBook return the indexes of routingBook.
func ListRoutingBook() (o []string) {
	for i := range routingBook {
		o = append(o, i)
	}
	sort.Strings(o)
	return
}
