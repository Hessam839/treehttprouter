package treehttprouter

import (
	"io"
	"net"
	"time"
)

type myConn struct {
	buff []byte
}

func (m *myConn) Read(b []byte) (n int, err error) {
	err = nil

	if len(b) == 0 {
		return 0, nil
	}

	n = copy(b, m.buff[:len(m.buff)])
	return
}

func (m *myConn) Write(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, io.EOF
	}
	l := len(m.buff)
	m.buff = append(m.buff, make([]byte, len(b))...)
	return copy(m.buff[l:], b), nil
}

func (m *myConn) Close() error {
	m.buff = m.buff[:0]
	return nil
}

func (m *myConn) LocalAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (m *myConn) RemoteAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (m *myConn) SetDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (m *myConn) SetReadDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (m *myConn) SetWriteDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func NewMockConn() *myConn {
	return &myConn{
		//buff: make([]byte, 1024),
	}
}
