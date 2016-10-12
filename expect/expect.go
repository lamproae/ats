package expect

import (
	"time"
	"log"
	"telnet"
)

const timeout = 10 * time.Second

func checkErr(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func Expect(c *telnet.Client, d ...string) {
	checkErr(c.SetReadDeadline(time.Now().Add(timeout)))
	checkErr(c.SkipUntil(d...))
}

func Sendln(c *telnet.Client, s string) {
	checkErr(c.SetWriteDeadline(time.Now().Add(timeout)))
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
	_, err := c.Write(buf)
	checkErr(err)
}
