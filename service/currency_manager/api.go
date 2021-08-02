package currency_manager

import (
	"errors"

	"github.com/NotSoFancyName/conversion_service/service/currency_manager/model"
	"github.com/NotSoFancyName/conversion_service/service/currency_manager/postgres"
)

type Manager interface {
	SaveCurrencies(entries []*model.CurrencyEntry) error
	GetCurrencies() ([]*model.CurrencyEntry, error)
	Shutdown() error
}

type ManagerType string

const (
	PostgresManager ManagerType = "postgres"
)

func NewManagerOfType(ttype ManagerType) (Manager, error) {
	switch ttype {
	case PostgresManager:
		return postgres.NewPostgreManager(), nil
	default:
		return nil, errors.New("unknown manager type")
	}
}
