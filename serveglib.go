//serveglib Help create server compatibly jsonrpc-glib-1.0
package jsonrpc2glib

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

var debug bool

func DebugMode(mode bool) {
	debug = mode
}
func NewMyConn(c net.Conn) *MyConn {
	p := &MyConn{Conn: c}
	p.Init()
	return p
}

type MyConn struct {
	Conn    net.Conn
	left    int
	reader1 *bufio.Reader
}

func (c *MyConn) Init() {
	c.reader1 = bufio.NewReader(c.Conn)
}

var parseErr = jsonrpc2.NewError(-32700, "parse error")

func (c *MyConn) Read(p []byte) (n int, err error) {
	if c.left == 0 {
		line1, _, err := c.reader1.ReadLine()
		if err != nil {
			return 0, err
		}
		if strings.HasPrefix(string(line1), "Content-Length:") == false {
			return 0, parseErr
		}
		lenData, err := strconv.ParseInt(string(line1[16:]), 10, 32)
		if err != nil {
			return 0, err
		}
		_, err = c.reader1.Discard(2)
		if err != nil {
			return 0, err
		}
		c.left = int(lenData)
	}
	lenP := len(p)
	if lenP >= c.left+1 {
		n, err = io.ReadFull(c.reader1, p[:c.left])
		if err != nil {
			return 0, err
		}
		p[c.left] = 10
		n = c.left + 1
		c.left = 0
		if debug {
			os.Stdout.Write(p[:n])
		}
		return n, nil
	} else {
		n, err = io.ReadFull(c.reader1, p[:lenP-1])
		if err != nil {
			return 0, err
		}
		c.left -= n
		if debug {
			os.Stdout.Write(p[:n])
		}
		return n, nil
	}
}

func (c *MyConn) Write(p []byte) (n int, err error) {
	if debug {
		os.Stdout.Write(p)
	}
	buf := bytes.NewBufferString("Content-Length: ")
	buf.WriteString(fmt.Sprintf("%d\r\n\r\n", len(p)-1))
	buf.Write(p[:len(p)-1])
	num, err := io.Copy(c.Conn, buf)
	return int(num), err
}

func (c *MyConn) Notify(method string, params interface{}) error {
	msgbuf := bytes.NewBufferString(`{"jsonrpc":"2.0","method":"`)
	msgbuf.WriteString(method)
	msgbuf.WriteString(`","params":`)
	param, err := json.Marshal(params)
	if err != nil {
		return err
	}
	msgbuf.Write(param)
	msgbuf.WriteString("}")
	_, err = c.Write(msgbuf.Bytes())
	return err
}

func (c *MyConn) Close() error {
	c.Conn.Close()
	return nil
}

//ServeGlib create serve, it will use DefaultServer if srv==nil.
func ServeGlib(conn net.Conn, srv *rpc.Server) {
	c := NewMyConn(conn)
	if srv == nil {
		srv = rpc.DefaultServer
	}
	rpc.ServeCodec(jsonrpc2.NewServerCodec(c, srv))
}
