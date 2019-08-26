package wss

import (
	"github.com/genshen/wssocks/wss/term_view"
	log "github.com/sirupsen/logrus"
	"net"
)

// listen on local address:port and forward socks5 requests to wssocks server.
func ListenAndServe(wsc *WebSocketClient, address string) error {
	s, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	log.WithField("local address", address).Info("listening on local address for incoming proxy request.")

	plog := term_view.NewPLog()
	log.SetOutput(plog) // change log stdout to plog

	var client Client
	for {
		c, err := s.Accept()
		if err != nil {
			return err
		}

		go func() {
			err := client.Reply(c, func(conn *net.TCPConn, proxyType int, addr string) error {
				proxy := wsc.NewProxy(conn)
				proxy.Serve(plog, wsc, proxyType, addr)
				wsc.TellClose(proxy.Id)
				return nil // todo error
			})
			if err != nil {
				log.Error(err)
			}
		}()
	}
}
