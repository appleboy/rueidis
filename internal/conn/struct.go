package conn

import (
	"bufio"
	proto2 "github.com/rueian/rueidis/internal/proto"
	"net"

	"github.com/rueian/rueidis/internal/queue"
)

type StructConn struct {
	r *bufio.Reader
	w *bufio.Writer
	q queue.Queue
}

func NewStructConn(conn net.Conn) *StructConn {
	c := &StructConn{r: bufio.NewReader(conn), w: bufio.NewWriter(conn), q: queue.NewRing(4096)}
	c.reading()
	return c
}

func (c *StructConn) reading() {
	go func() {
		for {
			t := c.q.Next1(true)
			if t == nil {
				if err := c.w.Flush(); err != nil {
					panic(err)
				}
				t = c.q.Next1(false)
			}
			if err := t.C.WriteTo(c.w); err != nil {
				panic(err)
			}
		}
	}()
	go func() {
		for {
			msg, err := proto2.ReadNextRaw(c.r)
			if err == nil {
				if msg.Type == '>' {
					continue
				}
			}
			t := c.q.Next2()
			t.W <- &resultS{S: msg, E: err}
		}
	}()
}

func (c *StructConn) Write(cmd proto2.StringArray) (proto2.Raw, error) {
	t := queue.Task{C: cmd, W: make(chan interface{}, 1)}
	c.q.Put(&t)
	r := (<-t.W).(*resultS)
	return r.S, r.E
}

type resultS struct {
	S proto2.Raw
	E error
}