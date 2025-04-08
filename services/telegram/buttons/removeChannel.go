package buttons

import (
	"awesomeProject/infrastructure/Postgres"
	"awesomeProject/services/twitch"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"strings"
)

func RemoveChannel(c tele.Context) error {
	brLogin := strings.TrimPrefix(c.Data()[1:], "rm.")

	if brLogin == "" {
		rows, err := getChannels(c.Sender().ID, "rm")
		if err != nil {
			return err
		}

		rows = append(rows, tele.Row{BtnBackToMenu})
		selector.Inline(rows...)

		return c.Edit("Выберите канал, который хотите удалить", selector)
	}

	err := twitch.RmChannel(brLogin, c.Sender().ID)
	selector.Inline(selector.Row(BtnBackToMenu))
	if err != nil {
		_ = c.Edit("Произошла внутренняя ошибка", selector)
		return err
	}
	return c.Edit("Канал удален", selector)

}

func getChannels(id int64, action string) (rows []tele.Row, err error) {
	pgx := Postgres.GetDB()

	channels, err := pgx.GetChannels(id)
	if err != nil {
		return nil, err
	}

	var (
		row []tele.Btn
	)
	for i, channel := range channels {
		row = append(row, tele.Btn{
			Text:   channel.BroadcasterLogin,
			Unique: fmt.Sprint(action, ".", channel.BroadcasterLogin),
		})

		if i%2 == 1 || i == len(channels)-1 {
			rows = append(rows, row)
			row = tele.Row{}
		}
	}
	return
}
