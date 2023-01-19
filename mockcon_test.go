package treehttprouter

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestNewMockConn(t *testing.T) {
	conn := NewMockConn()
	assert.NotNil(t, conn)
}

func TestGetLocalAddr(t *testing.T) {
	conn := NewMockConn()

	assert.EqualValues(t, conn.LocalAddr(), nil)
}

func TestRemoteAddr(t *testing.T) {
	conn := NewMockConn()

	assert.EqualValues(t, conn.RemoteAddr(), nil)
}

func TestRead(t *testing.T) {
	conn := NewMockConn()

	var buff []byte
	_, err := conn.Read(buff)
	assert.ErrorIs(t, err, io.ErrShortBuffer)
}

func TestWrite(t *testing.T) {
	conn := NewMockConn()

	var buff []byte
	_, err := conn.Write(buff)

	assert.ErrorIs(t, err, io.EOF)
}
