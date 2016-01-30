package testutil

import "github.com/google/go-gcm"

// TestCCS allow transparent testing of any functions depending on
// a CloudConnectionServer
type TestCCS struct {
	XMPPFunc      func(m *gcm.XmppMessage) (string, int, error)
	HTTPFunc      func(m *gcm.HttpMessage) (*gcm.HttpResponse, error)
	ListenMessage gcm.CcsMessage
}

// SendXMPP mocks an XMPP request by sending a message directly to the given
// testing function, XMPPFunc
func (ccs TestCCS) SendXMPP(m *gcm.XmppMessage) (string, int, error) {
	return ccs.XMPPFunc(m)
}

// SendHTTP mocks an HTTP request by sending a message directly to the given
// testing function, HTTPFunc
func (ccs TestCCS) SendHTTP(m *gcm.HttpMessage) (*gcm.HttpResponse, error) {
	return ccs.HTTPFunc(m)
}

// Listen mocks a CCS Listener by immediately sending the given testing message,
// ListenMessage to the given MessageHandler.
func (ccs TestCCS) Listen(h gcm.MessageHandler, stop <-chan bool) error {
	return h(ccs.ListenMessage)
}
