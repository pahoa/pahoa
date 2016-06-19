package core

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mattes/migrate/migrate"
)

type SQLModel struct {
	db *sql.DB
}

func NewSQLModel(driverName, dataSourceName, migrationPath string) (*SQLModel, error) {
	url := fmt.Sprintf("%s://%s", driverName, dataSourceName)
	log.Printf("Starting migration %s", url)
	errs, ok := migrate.UpSync(url, migrationPath)
	if !ok {
		return nil, errs[0]
	}
	log.Printf("Migration completed")

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return &SQLModel{db: db}, nil
}

func (m *SQLModel) GetCard(id string) (*Card, error) {
	log.Printf("GetCard(%s)", id)

	card := Card{}

	query := `
		select id, previous_step, current_step, status
		from card
		where id = ?
		`
	err := m.db.QueryRow(query, id).
		Scan(&card.ID, &card.PreviousStep, &card.CurrentStep, &card.Status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &card, nil
}

func (m *SQLModel) AddCard(opts *ModelAddCardOptions) (*Card, error) {
	log.Printf("AddCard(%#v)", opts)

	card, err := m.GetCard(opts.ID)
	if err != nil {
		return nil, err
	}
	if card != nil {
		return nil, ErrCardAlreadyExists
	}

	query := `
		insert into card(id, previous_step, current_step, status)
			values(?, ?, ?, ?)
		`
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(opts.ID, opts.PreviousStep, opts.CurrentStep, opts.Status)
	if err != nil {
		return nil, err
	}

	return m.GetCard(opts.ID)
}

func (m *SQLModel) UpdateCard(opts *ModelUpdateCardOptions) (*Card, error) {
	log.Printf("UpdateCard(%#v)", opts)

	card, err := m.GetCard(opts.ID)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, ErrCardNotFound
	}

	query := `
		update card set previous_step = ?, current_step = ?, status = ? 
			where id = ?
		`
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(opts.PreviousStep, opts.CurrentStep, opts.Status, opts.ID)
	if err != nil {
		return nil, err
	}

	return m.GetCard(opts.ID)
}

func (m *SQLModel) ListCards(step string) ([]*Card, error) {
	log.Printf("ListCards(%s)", step)

	var rows *sql.Rows
	var err error
	if step == "" {
		query := "select id, previous_step, current_step, status from card order by id"
		rows, err = m.db.Query(query)
		if err != nil {
			return nil, err
		}
	} else {
		query := `
			select id, previous_step, current_step, status 
				from card where current_step = ? 
				order by id
			`
		rows, err = m.db.Query(query, step)
		if err != nil {
			return nil, err
		}
	}

	defer rows.Close()

	cards := make([]*Card, 0)

	for rows.Next() {
		card := Card{}

		err := rows.Scan(&card.ID, &card.PreviousStep, &card.CurrentStep,
			&card.Status)
		if err != nil {
			return nil, err
		}
		cards = append(cards, &card)
	}

	return cards, nil
}

func (m *SQLModel) UpdateCardStatus(id, status string) error {
	log.Printf("UpdateCardStatus(%s, %s)", id, status)

	query := "update card set status = ? where id = ?"
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *SQLModel) ClearActionLogs(id string) error {
	log.Printf("ClearActionLogs(%s)", id)

	query := "delete from actionlog where card_id = ?"
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (m *SQLModel) CreateActionLog(id string, action Action) error {
	log.Printf("CreateActionLog(%s, %s)", id, action)

	query := `
		insert into actionlog(card_id, action, status, msg)
			values(?, ?, ?, ?)
		`
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, string(action), ActionStatusWaiting, "")
	if err != nil {
		return err
	}

	return nil
}

func (m *SQLModel) UpdateActionLogStatus(id string, action Action, status, msg string) error {
	log.Printf("UpdateActionLogStatus(%s, %s, %s, %s)", id, action, status, msg)

	query := "update actionlog set action = ?, status = ?, msg = ? where card_id = ?"
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(string(action), status, msg, id)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	migrate.NonGraceful()
}
