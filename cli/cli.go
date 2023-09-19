package cli

import (
	"flag"

	"github.com/google/uuid"
)

var AdminPass = flag.String("adminPass", uuid.New().String(), "admin password")
