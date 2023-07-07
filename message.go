package updaterproxy

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

type TargetClient struct {
	Vmuuid   string `json:"vmuuid"`
	Sn       string `json:"sn"`
	Hostname string `json:"hostname"`
	Ip       string `json:"ip"`
}

type Message struct {
	From     string          `json:"from"`
	To       string          `json:"to"`
	Id       string          `json:"id"`
	Type     string          `json:"type"`
	Method   string          `json:"method"`
	Data     json.RawMessage `json:"data"`
	Code     string          `json:"code"`
	Msg      string          `json:"msg"` // 新增 Msg 字段
	TraceId  string          `json:"traceId"`
	Timeout  time.Duration   // 添加 Timeout 字段
	ClientIP string          `json:"clientIp"`
	TaskId   string          `json:"taskId"`
}

const (
	METHOD_REQUEST  = "request"
	METHOD_RESPONSE = "response"

	CODE_SUCCESS = "success"
	CODE_ERROR   = "error"
	CODE_TIMEOUT = "timeout"
)

type Response struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type HandlerFunc func(ctx *Context) error

type MessageHandler struct {
	handlers map[string]HandlerFunc
	in       chan *Message
	out      chan *Message
}

func NewMessageHandler(bufferSize int) *MessageHandler {
	return &MessageHandler{
		handlers: make(map[string]HandlerFunc),
		in:       make(chan *Message, bufferSize),
		out:      make(chan *Message, bufferSize),
	}
}

func (h *MessageHandler) RegisterHandler(messageType string, handler HandlerFunc) {
	if _, exists := h.handlers[messageType]; exists {
		log.Fatalf("Handler already registered for message type: %s", messageType)
	}

	h.handlers[messageType] = handler
}

func (h *MessageHandler) HandleClientMessage(proxy *Proxy, num int) {
	for i := 0; i < num; i++ {
		go func() {
			// 用于防止 panic 造成的程序崩溃
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic in HandleMessages: %v", r)
				}
			}()

			for msg := range h.in {
				ctx0 := context.Background()
				if msg.Timeout > 0 {
					ctx0, _ = context.WithTimeout(ctx0, msg.Timeout)
				}

				handler, exists := h.handlers[msg.Type]
				if !exists {
					log.Printf("No handler registered for message type: %s", msg.Type)
					continue
				}

				ctx := &Context{
					Proxy:   proxy,
					Message: msg,
					Ctx:     ctx0,
				}

				if err := handler(ctx); err != nil {
					log.Printf("Error handling message: %v", err)
				}
			}
		}()
	}
}

func (h *MessageHandler) HandleServerMessage(proxy *Proxy, num int) {
	for i := 0; i < num; i++ {
		go func() {
			// 用于防止 panic 造成的程序崩溃
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic in HandleMessages: %v", r)
				}
			}()

			for msg := range h.out {
				ctx0 := context.Background()
				if msg.Timeout > 0 {
					ctx0, _ = context.WithTimeout(ctx0, msg.Timeout)
				}

				handler, exists := h.handlers[msg.Type]
				if !exists {
					log.Printf("No handler registered for message type: %s", msg.Type)
					continue
				}

				ctx := &Context{
					Proxy:   proxy,
					Message: msg,
					Ctx:     ctx0,
				}

				if err := handler(ctx); err != nil {
					log.Printf("Error handling message: %v", err)
				}
			}
		}()
	}
}
