package database

import (
	"avito/internal/pkg/config"
	"avito/internal/pkg/errors"
	"avito/internal/pkg/models"
	"database/sql"
	"fmt"
	"time"
)

type Pool struct {
	pool *sql.DB
}

func NewPool(c *config.DBConfig) (*Pool, error) {
	dsn := "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + c.Database

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(c.Connections)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Pool{pool: db}, nil
}

func (p *Pool) Close() error {
	return p.pool.Close()
}

type Repo struct {
	db *Pool
}

func NewRepo(DB *Pool) *Repo {
	return &Repo{
		db: DB,
	}
}

func (p *Repo) Transfer(input *models.Transfer) error {
	tr, err := p.db.pool.Begin()
	if err != nil {
		return errors.Wrap(err, "method Transfer: open transaction error")
	}
	defer tr.Rollback()

	var sender *uint = nil
	var recipient *uint = nil

	if input.RecipientID != 0 {
		recipient = &input.RecipientID
		result, err := tr.Exec(enrolmentQueryUpdate, input.Amount, input.RecipientID)
		if err != nil {
			return errors.Wrapf(err, "metod Transfer: update error, recipient_id: %d", input.RecipientID)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return errors.Wrapf(err, "metod Transfer: rows affected error, recipient_id: %d", input.RecipientID)
		} else if rows == 0 {
			_, err := tr.Exec(enrolmentQueryInsert, input.RecipientID, input.Amount)
			if err != nil {
				return errors.Wrapf(err, "metod Transfer: insert error, recipient_id: %d", input.RecipientID)
			}
		}
	}

	if input.SenderID != 0 {
		sender = &input.SenderID
		result, err := tr.Exec(writeOffQueryUpdate, input.Amount, input.SenderID)
		if err != nil {
			return errors.InsufficientFunds.Wrapf(err, "method Transfer: update balance error, sender_id: %d", input.SenderID)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return errors.Wrapf(err, "metod Transfer: rows affected error, sender_id: %d", input.SenderID)
		} else if rows == 0 {
			return errors.InsufficientFunds.Newf("method Transfer: update balance error, sender_id: %d", input.SenderID)
		}
	}

	_, err = tr.Exec(transferQueryInsert,
		sender, recipient, input.Amount, time.Now().Format(time.RFC3339), input.Message)
	if err != nil {
		return errors.Wrap(err, "method Transfer: insert in transfers error")
	}

	err = tr.Commit()
	if err != nil {
		return errors.Wrap(err, "method Transfer: commit error")
	}

	return nil
}

func (p *Repo)  GetBalance(id uint) (*models.Balance, error) {
	balance := &models.Balance{}

	err := p.db.pool.QueryRow(balanceQuerySelect,
		id,
	).Scan(&balance.UserID, &balance.Total)
	if err != nil {
		return nil, errors.BalanceNotFound.Wrap(err, "method GetBalance: query error")
	}

	return balance, nil
}

func (p *Repo) GetListTransfer(id uint, sort string, order string, limit int, offset int) (models.TransferList, error) {
	list:= make(models.TransferList, 0)

	rows, err := p.db.pool.Query(fmt.Sprintf(transferQuerySelect, sort, order), id, limit, offset)
	if err != nil {
		return nil, errors.Wrapf(err, "method GetListTransfer: query error, user_id: %d", id)
	}
	defer rows.Close()

	for rows.Next() {
		dbTransfer := &models.DBTransfer{}
		err = rows.Scan(&dbTransfer.SenderID, &dbTransfer.RecipientID, &dbTransfer.Amount, &dbTransfer.Date, &dbTransfer.Message)
		if err != nil {
			return nil, errors.Wrapf(err, "method GetListTransfer: scaning rows error, user_id: %d", id)
		}

		list = append(list, dbTransfer.GetPin())
	}

	return list, nil
}

const (
	writeOffQueryUpdate = "UPDATE users_balance SET balance = balance - $1 WHERE user_id = $2"
	enrolmentQueryUpdate = "UPDATE users_balance SET balance = balance + $1 WHERE user_id = $2"
	enrolmentQueryInsert = "INSERT INTO users_balance (user_id, balance) VALUES($1, $2)"
	transferQueryInsert = "INSERT INTO transfers (sender_id, recipient_id, amount, date, message) VALUES($1, $2, $3, $4, $5)"
	balanceQuerySelect = "SELECT user_id, balance FROM users_balance WHERE user_id = $1"
	transferQuerySelect = "SELECT sender_id, recipient_id, amount, date, message FROM transfers" +
		" WHERE sender_id = $1 or recipient_id = $1 ORDER BY %s %s" +
		" LIMIT $2 OFFSET $3"
)
