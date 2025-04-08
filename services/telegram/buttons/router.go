package buttons

import (
	tele "gopkg.in/telebot.v3"
	"strings"
)

func Router(c tele.Context) error {
	err := c.Respond()
	if err != nil {
		return err
	}
	data := strings.TrimPrefix(c.Data(), "\f")
	switch {
	case strings.HasPrefix(data, "rm"):
		return RemoveChannel(c)
	case strings.HasPrefix(data, "add"):
		return AddChannel(c)
	case strings.HasPrefix(data, "aT"):
		return AddTrigger(c)
	case strings.HasPrefix(data, "rT"):
		return RemoveTrigger(c)
	case strings.HasPrefix(data, "s."):
		return Silence(c)
	case strings.HasPrefix(data, "uS."):
		return UnSilence(c)
	case data == "cancel":
		mutex.Lock()
		defer mutex.Unlock()
		userStates[c.Sender().ID] = 0
		return MenuCommand(c)
	}
	return nil
}
