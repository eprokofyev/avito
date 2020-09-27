package database

import "avito/internal/pkg/models"

type IDatabase interface{
	Transfer(*models.Transfer) error
	GetBalance(uint) (*models.Balance, error)
	GetListTransfer(id uint, sort string, order string, limit int, offset int) (models.TransferList, error)
}
