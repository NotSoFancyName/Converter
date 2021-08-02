package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/NotSoFancyName/conversion_service/service/currency_manager"
	"github.com/NotSoFancyName/conversion_service/service/currency_manager/model"
)

const (
	currencyURL = "http://api.currencylayer.com/live"
	apiKey      = "ac1dcc2f8f5173e6d44771e0a7e9bc8f"

	baseCurrency = "USD"
)

type Fetcher struct {
	fetchPeriod time.Duration
	client      *http.Client
	url         *url.URL
	cm          currency_manager.Manager
}

func NewFetcher(period time.Duration) (*Fetcher, error) {
	u, err := url.Parse(currencyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}
	q := u.Query()
	q.Set("access_key", apiKey)
	q.Set("currencies", model.AllCurrenciesTypesString())
	u.RawQuery = q.Encode()

	cm, err := currency_manager.NewManagerOfType(currency_manager.PostgresManager)
	if err != nil {
		return nil, err
	}

	return &Fetcher{
		fetchPeriod: period,
		client:      &http.Client{},
		url:         u,
		cm:          cm,
	}, nil
}

func (f *Fetcher) Run(stop chan struct{}) {
	ticker := time.NewTicker(f.fetchPeriod)
	defer ticker.Stop()

	log.Println("Making initial fetch")
	entries, err := f.fetch()
	if err != nil {
		log.Printf("failed to fetch exchange rates: %v\n", err)
	}
	err = f.cm.SaveCurrencies(entries)
	if err != nil {
		log.Printf("failed to save exchange rates: %v\n", err)
	}

	for {
		select {
		case <-stop:
			log.Println("Stopping fetcher service")
			stop <- struct{}{}
			return
		case <-ticker.C:
			log.Println("Trying to fetch currency exachange rates")
			entries, err = f.fetch()
			if err != nil {
				log.Printf("failed to fetch exchange rates: %v\n", err)
			}
			err = f.cm.SaveCurrencies(entries)
			if err != nil {
				log.Printf("failed to save exchange rates: %v\n", err)
			}
		}
	}
}

func (f *Fetcher) fetch() ([]*model.CurrencyEntry, error) {
	resp, err := f.client.Get(f.url.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}
	defer resp.Body.Close()

	var currencyRatios struct {
		Status bool               `json:"success"`
		Source string             `json:"source"`
		Quotes map[string]float64 `json:"quotes"`
	}

	err = json.NewDecoder(resp.Body).Decode(&currencyRatios)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body': %v", err)
	}

	if !currencyRatios.Status {
		return nil, errors.New("the success status is false")
	}

	var entries []*model.CurrencyEntry
	for k, v := range currencyRatios.Quotes {
		cur := strings.TrimPrefix(k, baseCurrency)
		entries = append(entries,
			&model.CurrencyEntry{
				Type:  model.CurrencyType(cur),
				Ratio: v,
			},
		)
	}

	return entries, nil
}
