package bot

import (
	"context"

	"github.com/go-telegram/bot/models"
)

func applyMiddlewares(h HandlerFunc, m ...Middleware) HandlerFunc {
	if len(m) < 1 {
		return h
	}
	wrapped := h
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}
	return wrapped
}

// ProcessUpdate allows you to process update
func (b *Bot) ProcessUpdate(ctx context.Context, upd *models.Update) {
	h := b.defaultHandlerFunc

	if upd.Message != nil {
		h = b.findHandler(HandlerTypeMessageText, upd)
	} else if upd.CallbackQuery != nil {
		h = b.findHandler(HandlerTypeCallbackQueryData, upd)
	}

	applyMiddlewares(h, b.middlewares...)(ctx, b, upd)
}

func (b *Bot) findHandler(handlerType HandlerType, upd *models.Update) HandlerFunc {
	b.handlersMx.RLock()
	defer b.handlersMx.RUnlock()

	for _, h := range b.handlers {
		if h.handlerType == handlerType {
			if h.match(upd) {
				return h.handler
			}
		}
	}

	return b.defaultHandlerFunc
}
