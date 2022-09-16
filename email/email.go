package email

import (
	"log"
	"time"

	"gopkg.in/gomail.v2"
)

type Client struct {
	dialer     *gomail.Dialer
	sendCloser gomail.SendCloser
	ch         chan *gomail.Message
	open       bool
}

func New(host string, port int, username string, password string) *Client {
	c := &Client{
		dialer: gomail.NewDialer(host, port, username, password),
		ch:     make(chan *gomail.Message),
	}
	go c.run()
	return c
}

func (c *Client) Send(m *gomail.Message) {
	c.ch <- m
}

func (c *Client) run() {
	var err error
	timer := time.NewTimer(0)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()
	for {
		// Reuse timer
		if !timer.Stop() {
			<-timer.C
		}
		timer.Reset(30 * time.Second)

		select {
		case m, ok := <-c.ch:
			if !ok {
				continue
			}
			if !c.open {
				if c.sendCloser, err = c.dialer.Dial(); err != nil {
					continue
				}
				c.open = true
			}
			if err := gomail.Send(c.sendCloser, m); err != nil {
				log.Print(err)
			}
		case <-timer.C:
			if c.open {
				if err := c.sendCloser.Close(); err != nil {
					continue
				}
				c.open = false
			}
		}
	}
}
