package cli

import (
	"errors"
	"net/http"
	"os"
)

// Thanks to https://stackoverflow.com/a/57281956 for this idea and code.

// wrappedDir wraps http.Dir so that we can modify http.Dir's Open method.
type WrappedDir struct {
	dir http.Dir
}

// Open modifies http.Dir's Open method in order to handle URLs without ".html".
func (wd WrappedDir) Open(name string) (http.File, error) {
	fh, err := wd.dir.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		// errPlus == nil because I'm (unusually) looking for a successful open.
		if fhPlus, errPlus := wd.dir.Open(name + ".html"); errPlus == nil {
			return fhPlus, nil
		}
	}
	return fh, err
}
