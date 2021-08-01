package mods

import (
	"os"
	"path/filepath"
)

var PgLoc, _ = os.Executable()
var PgPath = filepath.Dir(PgLoc)
