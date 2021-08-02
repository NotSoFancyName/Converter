package postgres

import (
	"errors"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/NotSoFancyName/conversion_service/service/currency_manager/model"
)

const (
	reconnectTime = 3 * time.Second
)

type PostgreManager struct {
	db *gorm.DB
}

func NewPostgreManager() *PostgreManager {
	for {
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "host=postgres user=postgres password=pwadmin dbname=converter-db port=5432",
			PreferSimpleProtocol: true,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Printf("Failed to connect to DB: %v", err)
			time.Sleep(reconnectTime)
			continue
		}
		db.AutoMigrate(&Currency{})
		return &PostgreManager{
			db: db,
		}
	}
}

type Currency struct {
	gorm.Model
	Type  string
	Ratio float64
}

func (m *PostgreManager) GetCurrencies() ([]*model.CurrencyEntry, error) {
	var cur []Currency
	m.db.Limit(len(model.AllCurrenciesTypes)).Order("created_at desc").Find(&cur)
	if len(cur) == 0 {
		return nil, errors.New("the database does not contain currency rates yet")
	}

	var res []*model.CurrencyEntry
	for _, v := range cur {
		res = append(res, &model.CurrencyEntry{
			Type:  model.CurrencyType(v.Type),
			Ratio: v.Ratio,
		})
	}
	return res, nil
}

func (m *PostgreManager) SaveCurrencies(entries []*model.CurrencyEntry) error {
	var currencies []Currency
	for _, v := range entries {
		currencies = append(currencies,
			Currency{
				Type:  string(v.Type),
				Ratio: v.Ratio,
			},
		)
	}
	m.db.Create(&currencies)
	return nil
}

func (m *PostgreManager) Shutdown() error {
	db, err := m.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
