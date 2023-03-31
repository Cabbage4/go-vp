package butin

import (
	"strings"
)

var (
	skipError = []string{"connection reset by peer", "use of closed network connection", "EOF"}
)

func IsSkipError(err error) bool {
	for _, v := range skipError {
		if strings.Contains(err.Error(), v) {
			return true
		}
	}

	return false
}
