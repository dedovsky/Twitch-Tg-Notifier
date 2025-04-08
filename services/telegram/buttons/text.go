package buttons

import (
	"awesomeProject/services/twitch"
	tele "gopkg.in/telebot.v3"
	"strings"
	"sync"
)

const (
	AddChannelState = 1 + iota
	AddTriggerName
)

var (
	userStates = make(map[int64]int)

	mutex = sync.RWMutex{}
)

var (
	selector = &tele.ReplyMarkup{}

	BtnBackToMenu = selector.Data("🏠 Вернуться в меню", "menu")

	BtnAddChannel    = selector.Data("Добавить новый канал", "add")
	BtnRemoveChannel = selector.Data("Удалить канал", "rm.")

	BtnNotifySettings = selector.Data("Настройки уведомлений", "notifySettings")
)

func MenuCommand(c tele.Context) error {
	selector.Inline(selector.Row(BtnAddChannel, BtnRemoveChannel), selector.Row(BtnNotifySettings))

	text := "test"

	if c.Text() != "/start" {
		return c.Edit(text, selector)
	}
	return c.Send("тест", selector)
}

func OnText(c tele.Context) error {
	mutex.RLock()
	defer mutex.RUnlock()

	switch userStates[c.Sender().ID] {
	case AddChannelState:
		delete(userStates, c.Sender().ID)
		return addChannelState(c)
	case AddTriggerName:
		delete(userStates, c.Sender().ID)
		return addTrigger(c)
	}
	return nil
}

func addChannelState(c tele.Context) error {
	var brLogin string
	var found bool
	if !strings.HasPrefix(c.Text(), "https://www.twitch.tv") && !strings.HasPrefix(c.Text(), "https://twitch.tv") {
		brLogin = c.Text()
	} else {
		brLogin, found = strings.CutPrefix(c.Text(), "https://www.twitch.tv/")
		if !found {
			brLogin, found = strings.CutPrefix(c.Text(), "https://twitch.tv/")
			if !found {
				return c.Send("Некорректная ссылка на канал")
			}
		}
	}
	selector.Inline(selector.Row(BtnAddChannel, BtnRemoveChannel), selector.Row(BtnNotifySettings))
	rowsAffected, err := twitch.AddChannel(brLogin, c.Sender().ID)
	if err != nil {
		_ = c.Send("Произошла внутренняя ошибка", selector)
		return err
	}
	if rowsAffected == 0 {
		return c.Send("Канал уже добавлен", selector)
	}
	return c.Send("Канал добавлен", selector)
}
