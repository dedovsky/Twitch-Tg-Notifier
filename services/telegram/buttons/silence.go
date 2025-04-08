package buttons

import (
	"awesomeProject/services/twitch"
	tele "gopkg.in/telebot.v3"
	"strings"
)

func Silence(c tele.Context) error {
	brLogin := strings.TrimPrefix(c.Data(), "\fs.")

	if brLogin == "" {
		return c.RespondAlert("Произошла ошибка")
	}

	selector.Inline(selector.Row(selector.Data("🔕", "uS."+brLogin)))
	_ = c.Edit(c.Text(), selector)

	err := twitch.SetSilenced(brLogin, c.Sender().ID, true)
	return err
}

func UnSilence(c tele.Context) error {
	brLogin := strings.TrimPrefix(c.Data(), "\fuS.")

	if brLogin == "" {
		return c.RespondAlert("Произошла ошибка")
	}

	selector.Inline(selector.Row(selector.Data("🔔", "s."+brLogin)))
	_ = c.Edit(c.Text(), selector)

	err := twitch.SetSilenced(brLogin, c.Sender().ID, false)
	return err
}
