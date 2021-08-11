package commands

import tg "gopkg.in/tucnak/telebot.v2"

type Command interface {
	Execute(m *tg.Message)
}
