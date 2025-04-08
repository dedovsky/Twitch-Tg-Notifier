package Postgres

import "github.com/jackc/pgx/v5"

type Trigger struct {
	Id      int64
	UserID  int64
	Trigger string
}

func (p *Postgres) AddTrigger(userID int64, trigger string) error {
	_, err := p.conn.Exec(p.ctx, `INSERT INTO triggers (user_id, trigger) VALUES ($1, $2)`, userID, trigger)
	return err
}

func (p *Postgres) RmTrigger(id string) error {
	_, err := p.conn.Exec(p.ctx, `DELETE FROM triggers WHERE id = $1`, id)
	return err
}

func (p *Postgres) GetTriggersForUser(userID int64) ([]Trigger, error) {
	var rows pgx.Rows
	var err error

	rows, err = pgxStruct.conn.Query(pgxStruct.ctx, "SELECT * FROM triggers WHERE user_id = $1", userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var triggers []Trigger

	for rows.Next() {
		tr, err := pgx.RowToStructByName[Trigger](rows)
		if err != nil {
			return nil, err
		}
		triggers = append(triggers, tr)
	}

	return triggers, nil
}

func (p *Postgres) GetTriggers() ([]string, error) {
	var rows pgx.Rows
	var err error

	rows, err = p.conn.Query(p.ctx, "SELECT trigger FROM triggers")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var triggers []string

	for rows.Next() {
		var tr string
		err = rows.Scan(&tr)
		if err != nil {
			return nil, err
		}
		triggers = append(triggers, tr)
	}

	return triggers, nil
}

func (p *Postgres) GetTriggerAuthors(trigger string) ([]int64, error) {
	var rows pgx.Rows
	var err error

	rows, err = p.conn.Query(p.ctx, "SELECT user_id FROM triggers WHERE trigger = $1", trigger)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var authors []int64

	for rows.Next() {
		var author int64
		err = rows.Scan(&author)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}
