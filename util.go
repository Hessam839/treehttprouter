package treehttprouter

import "strings"

func split(path string) (string, string) {
	if len(path) == 0 {
		return "/", ""
	}
	p := strings.Split(path, "/")

	if p[0] == "" {
		return "/", strings.Join(p[1:], "/")
	}
	return p[0], strings.Join(p[1:], "/")
}
