package models

import (
	"database/sql"
	"time"
)

//easyjson:json
type TransferList []*Transfer

//easyjson:json
type Transfer struct {
	SenderID uint `json:"sender_id,omitempty"`
	RecipientID uint `json:"recipient_id,omitempty"`
	Amount float64 `json:"amount"`
	Message string `json:"message,omitempty"`
	Date time.Time `json:"date,omitempty"`
}

type DBTransfer struct {
	SenderID sql.NullInt32
	RecipientID sql.NullInt32
	Amount float64
	Message sql.NullString
	Date sql.NullTime
}

func (dbt *DBTransfer) GetPin() *Transfer {
	tmp := &Transfer {
		Amount: dbt.Amount,
	}

	if dbt.SenderID.Valid {
		tmp.SenderID = uint(dbt.SenderID.Int32)
	} else {
		tmp.SenderID = 0
	}

	if dbt.RecipientID.Valid {
		tmp.RecipientID = uint(dbt.RecipientID.Int32)
	} else {
		tmp.RecipientID = 0
	}

	if dbt.Message.Valid {
		tmp.Message = dbt.Message.String
	} else {
		tmp.Message = ""
	}

	if dbt.Date.Valid {
		tmp.Date = dbt.Date.Time
	} else {
		tmp.Date = time.Time{}
	}

	return tmp
}
