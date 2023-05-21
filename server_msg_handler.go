package updaterproxy

type ServerMsgHandler struct {
}

func (handler *ServerMsgHandler) Request(ctx *Context) error {
	ctx.Proxy.SendToClient(ctx.Message)
	return nil
}
