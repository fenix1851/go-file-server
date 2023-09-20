package cli

import (
	"flag"

	"github.com/google/uuid"
)

var DefaultAdminPass = uuid.New().String()
var AdminPass = flag.String("adminPass", DefaultAdminPass, "admin password")
