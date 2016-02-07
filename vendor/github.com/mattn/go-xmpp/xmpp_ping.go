package xmpp

import (
	"fmt"
)

func (c *Client) PingC2S(jid, server string) error {
	_, err := fmt.Fprintf(c.conn, "<iq from='%s' to='%s' id='c2s1' type='get'>\n"+
		"<ping xmlns='urn:xmpp:ping'/>\n"+
		"</iq>",
		xmlEscape(jid), xmlEscape(server))
	return err
}

func (c *Client) PingS2S(fromServer, toServer string) error {
	_, err := fmt.Fprintf(c.conn, "<iq from='%s' to='%s' id='s2s1' type='get'>\n"+
		"<ping xmlns='urn:xmpp:ping'/>\n"+
		"</iq>",
		xmlEscape(fromServer), xmlEscape(toServer))
	return err
}