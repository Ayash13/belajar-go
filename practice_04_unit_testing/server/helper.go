package server

import "fmt"

func NewAPIPath(method string, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}
