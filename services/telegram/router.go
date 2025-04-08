package telegram

import (
	"awesomeProject/services/telegram/buttons"
	tele "gopkg.in/telebot.v3"
)

func PushRoutes(b *tele.Bot) {
	b.Handle("/start", buttons.MenuCommand)
	b.Handle(tele.OnText, buttons.OnText)

	b.Handle(tele.OnCallback, buttons.Router)
	b.Handle(&buttons.BtnBackToMenu, buttons.MenuCommand)
	b.Handle(&buttons.BtnNotifySettings, buttons.NotifySettings)
}
