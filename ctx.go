package updaterproxy

import "context"

type Context struct {
	Proxy   *Proxy
	Message *Message
	Extra   map[string]interface{}
	Ctx     context.Context
	Cancel  context.CancelFunc
}
