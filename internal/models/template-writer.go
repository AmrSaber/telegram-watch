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

func (w *TelegramTemplateWriter) Write(bytes []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.content = append(w.content, bytes...)

	if err := w.flush(); err != nil {
		return 0, err
	}

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
func (w *TelegramTemplateWriter) flush() error {
	msgContent := fmt.Appendf([]byte{}, w.template, w.content)

	if err := w.writer.SetContent(msgContent); err != nil {
		// FIXME: this error is lost, the command at runner.go sends "signal: broken pipe" error
		return err
	}

	w.started = true
	return nil
}
