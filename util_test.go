package treehttprouter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitPath(t *testing.T) {
	p1 := "/api/v2/users"
	p2 := ""
	expected := []string{"/", "api", "v2", "users"}
	i := 0
	for {
		p1, p2 = split(p1)
		assert.Equal(t, p1, expected[i])
		if p2 == "" {
			break
		}
		p1 = p2
		p2 = ""
		i++
	}
}

func TestSplitEmptyPath(t *testing.T) {
	p1, _ := split("")

	assert.Equal(t, p1, "/")
}
