package cli

import (
	"flag"
)

var AdminPass = flag.String("adminPass", "", "admin password")
var Port = flag.String("port", ":4001", "port to run server on")
