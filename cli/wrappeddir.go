package cli

import (
	"errors"
	"net/http"
	"os"
)

// Code care of https://stackoverflow.com/a/57281956/26702.

// wrappedDir wraps http.Dir so that we can wrap http.Dir's Open method.
type wrappedDir struct {
	d http.Dir
}

// Open wraps http.Dir's Open method in order to handle URLs without ".html".
func (d wrappedDir) Open(name string) (http.File, error) {
	f, err := d.d.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		if fPlus, errPlus := d.d.Open(name + ".html"); errPlus == nil {
			return fPlus, nil
		}
	}

	return f, err
}
