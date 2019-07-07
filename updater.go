package main

import (
	"regexp"
	"strings"
	"time"
)

func update(pt *string) updater {
	key := false
	return func(term goraptor.Term) error {
		if set {
			return fmt.Errorf("Property Already Defined")
		}

		pt = termStr(term)
		key = true
		return nil
	}
}