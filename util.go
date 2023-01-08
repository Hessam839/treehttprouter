package treehttprouter

import "strings"

func split(path string) (string, string) {
	p := strings.Split(path, "/")
	if len(p) == 0 {
		return "/", ""
	}
	if p[0] == "" {
		return "/", strings.Join(p[1:], "/")
	}
	return p[0], strings.Join(p[1:], "/")
}
