package twitch

import (
	"awesomeProject/infrastructure/Postgres"
	"fmt"
	"github.com/gempir/go-twitch-irc/v4"
	tele "gopkg.in/telebot.v3"
	"log"
	"strings"
)

var bot = twitch.NewAnonymousClient()

var triggers []string

func ListenAndServe(b *tele.Bot) {
	db := Postgres.GetDB()

	channels, err := db.GetChannels()
	if err != nil {
		log.Println("Ошибка при получении каналов:", err)
		return
	}

	triggers, err = db.GetTriggers()
	if err != nil {
		log.Println("Ошибка при получении триггеров:", err)
		return
	}

	for _, channel := range channels {
		bot.Join(channel.BroadcasterLogin)
	}

	bot.OnPrivateMessage(func(message twitch.PrivateMessage) {
		msgLower := strings.ToLower(message.Message)

		for _, trigger := range triggers {
			if strings.Contains(msgLower, trigger) {
				authors, err := Postgres.GetDB().GetTriggerAuthors(trigger)
				if err != nil {
					log.Println("Ошибка при получении авторов триггера:", err)
					continue
				}
				for _, author := range authors {
					selector := tele.ReplyMarkup{}
					silenced := Postgres.GetDB().IsSilenced(author, message.Channel)

					var (
						icon   string
						action string
						silent bool
					)

					if silenced {
						icon = "🔕"
						action = "uS." + message.Channel // unSilence
						silent = true
					} else {
						icon = "🔔"
						action = "s." + message.Channel // silence
						silent = false
					}

					selector.Inline(selector.Row(selector.Data(icon, action)))

					_, err := b.Send(
						tele.ChatID(author),
						fmt.Sprintf("В #%s\n\n<b><u>%s</u></b>: %s", message.Channel, message.User.DisplayName, message.Message),
						&tele.SendOptions{DisableNotification: silent, ReplyMarkup: &selector, ParseMode: tele.ModeHTML},
					)

					if err != nil {
						log.Println("Ошибка при отправке сообщения:", err)
					}
				}
			}
		}
	})

	if err := bot.Connect(); err != nil {
		log.Println("Ошибка при подключении к Twitch:", err)
	}
}

func AddTrigger(userID int64, trigger string) error {
	err := Postgres.GetDB().AddTrigger(userID, trigger)
	if err != nil {
		return err
	}
	triggers = append(triggers, trigger)
	return nil
}

func RmTrigger(trigger string, trID string) error {
	err := Postgres.GetDB().RmTrigger(trID)
	if err != nil {
		return err
	}
	for i, t := range triggers {
		if t == trigger {
			triggers = append(triggers[:i], triggers[i+1:]...)
			break
		}
	}
	return nil
}

func AddChannel(brLogin string, userID int64) (int64, error) {
	bot.Join(brLogin)
	return Postgres.GetDB().AddChannel(brLogin, userID)
}

func RmChannel(brLogin string, userID int64) error {
	bot.Depart(brLogin)
	return Postgres.GetDB().RmChannel(brLogin, userID)
}

func SetSilenced(brLogin string, userID int64, silence bool) error {
	return Postgres.GetDB().SetSilenced(userID, brLogin, silence)
}
