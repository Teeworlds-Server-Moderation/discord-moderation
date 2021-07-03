package service

import (
	"fmt"
)

func fmtError(err error) string {
	return fmt.Sprintf("[ERROR]: %s", err)
}
