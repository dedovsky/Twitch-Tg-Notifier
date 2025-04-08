package buttons

import (
	"awesomeProject/infrastructure/Postgres"
	"awesomeProject/services/twitch"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"strings"
)

func NotifySettings(c tele.Context) error {
	selector.Inline(selector.Row(
		selector.Data("Добавить триггер", "aT"),
		selector.Data("Удалить триггер", "rT.")),
		selector.Row(BtnBackToMenu))

	return c.Edit("Настройки уведомлений", selector)
}

func AddTrigger(c tele.Context) error {
	selector.Inline(selector.Row(selector.Data("Отмена", "cancel")))

	userStates[c.Sender().ID] = AddTriggerName
	return c.Edit("Введите триггер", selector)

}

func addTrigger(c tele.Context) error {
	err := twitch.AddTrigger(c.Sender().ID, strings.ToLower(c.Text()))
	if err != nil {
		return err
	}

	selector.Inline(selector.Row(BtnBackToMenu))
	return c.Send("Триггер добавлен", selector)
}

func RemoveTrigger(c tele.Context) error {
	data := strings.TrimPrefix(c.Data(), "\f")
	trID := strings.TrimPrefix(data, "rT.")

	if trID == "" {
		rows, err := getTriggers(c.Sender().ID, "rT")
		if err != nil {
			return err
		}

		rows = append(rows, tele.Row{BtnBackToMenu})
		selector.Inline(rows...)

		return c.Edit("Выберите триггер, который хотите удалить", selector)
	}

	err := twitch.RmTrigger(c.Text(), trID)
	selector.Inline(selector.Row(BtnBackToMenu))
	if err != nil {
		_ = c.Edit("Произошла внутренняя ошибка", selector)
		return err
	}
	return c.Edit("Триггер удален", selector)
}

func getTriggers(id int64, action string) (rows []tele.Row, err error) {
	pgx := Postgres.GetDB()

	triggers, err := pgx.GetTriggersForUser(id)
	if err != nil {
		return nil, err
	}

	var (
		row []tele.Btn
	)
	for i, trigger := range triggers {
		row = append(row, tele.Btn{
			Text:   trigger.Trigger,
			Unique: fmt.Sprint(action, ".", trigger.Id),
		})

		if i%2 == 1 || i == len(triggers)-1 {
			rows = append(rows, row)
			row = []tele.Btn{}
		}
	}
	return
}
