package models

import (
	"fmt"
	"sync"
)

type TelegramTemplateWriter struct {
	lock sync.Mutex

	writer *TelegramWriter

	started bool

	content  []byte
	template string // Must contain exactly one "%s" to be replaces with the content
}

func NewTelegramTemplateWriter(userId string, template string) *TelegramTemplateWriter {
	writer := NewTelegramWriter(userId)

	return &TelegramTemplateWriter{
		writer: writer,

		content:  make([]byte, 0),
		template: template,
	}
}

func (w *TelegramTemplateWriter) GetChatId() int64 {
	return w.writer.GetChatId()
}

func (w *TelegramTemplateWriter) Wait() {
	w.writer.Wait()
}

func (w *TelegramTemplateWriter) Write(bytes []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.content = append(w.content, bytes...)
	w.flush()

	return len(bytes), nil
}

func (w *TelegramTemplateWriter) SetTemplate(template string) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.template = template
	if w.started {
		w.flush()
	}
}

// Write latest content to message
func (w *TelegramTemplateWriter) flush() {
	msgContent := fmt.Appendf([]byte{}, w.template, w.content)
	w.writer.SetContent(msgContent)
	w.started = true
}
