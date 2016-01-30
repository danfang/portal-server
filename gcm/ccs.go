package main

import "github.com/google/go-gcm"

// A CloudConnectionServer sends and receives GCM messages.
type CloudConnectionServer interface {
	SendXMPP(m *gcm.XmppMessage) (string, int, error)
	SendHTTP(m *gcm.HttpMessage) (*gcm.HttpResponse, error)
	Listen(h gcm.MessageHandler, stop <-chan bool) error
}

// GoogleCCS is the real GCM server, which will authenticate
// with a given Sender ID (project number) and API key.
type GoogleCCS struct {
	SenderID string
	APIKey   string
}

// SendXMPP sends an XMPP message via Google's GCM service
func (ccs GoogleCCS) SendXMPP(m *gcm.XmppMessage) (string, int, error) {
	return gcm.SendXmpp(ccs.SenderID, ccs.APIKey, *m)
}

// SendHTTP sends an HTTP message via Google's GCM service
func (ccs GoogleCCS) SendHTTP(m *gcm.HttpMessage) (*gcm.HttpResponse, error) {
	return gcm.SendHttp(ccs.APIKey, *m)
}

// Listen receives incoming GCM messages via Google's GCM service
func (ccs GoogleCCS) Listen(h gcm.MessageHandler, stop <-chan bool) error {
	return gcm.Listen(ccs.SenderID, ccs.APIKey, h, stop)
}
