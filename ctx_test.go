package treehttprouter

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCtx(t *testing.T) {
	//client := NewMockConn()
	req, _ := http.NewRequest("Get", "/login", bytes.NewReader([]byte(`{"name":"hesam","age":42}`)))

	//var buff bytes.Buffer
	//_ = req.Write(&buff)
	//_, _ = client.Write(buff.Bytes())

	ctx, err := NewCtx(req)
	if err != nil {
		t.Fatalf("contextect creation failed: [%v]", err)
	}
	assert.NotNil(t, ctx)

}
