package Postgres

import (
	"github.com/jackc/pgx/v5"
	"log"
)

type Channel struct {
	Id               int64
	UserID           int64
	BroadcasterLogin string
	Silence          bool
}

func (p *Postgres) AddChannel(BrLogin string, userID int64) (int64, error) {
	tag, err := p.conn.Exec(p.ctx, `INSERT INTO channels (user_id, broadcaster_login) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, BrLogin)

	return tag.RowsAffected(), err
}

func (p *Postgres) RmChannel(brLogin string, userID int64) error {
	_, err := p.conn.Exec(p.ctx, `DELETE FROM channels WHERE user_id = $1 and broadcaster_login = $2`, userID, brLogin)
	return err
}

func (p *Postgres) IsSilenced(userID int64, brLogin string) bool {
	var silenced bool
	err := p.conn.QueryRow(p.ctx, `SELECT silence FROM channels WHERE user_id = $1 AND broadcaster_login = $2`, userID, brLogin).Scan(&silenced)
	if err != nil {
		log.Println("Error checking silenced status:", err)
	}
	return silenced
}

func (p *Postgres) SetSilenced(userID int64, brLogin string, silence bool) error {
	_, err := p.conn.Exec(p.ctx, `UPDATE channels SET silence = $1 WHERE user_id = $2 AND broadcaster_login = $3`, silence, userID, brLogin)
	return err
}

func (p *Postgres) GetChannels(userID ...int64) ([]Channel, error) {
	var rows pgx.Rows
	var err error
	if len(userID) > 0 {
		rows, err = p.conn.Query(p.ctx, "SELECT * FROM channels WHERE user_id = $1", userID[0])
	} else {
		rows, err = p.conn.Query(p.ctx, "SELECT * FROM channels")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var channels []Channel

	for rows.Next() {
		ch, err := pgx.RowToStructByName[Channel](rows)
		if err != nil {
			return nil, err
		}
		channels = append(channels, ch)
	}

	return channels, nil
}
