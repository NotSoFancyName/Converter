package rest

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/NotSoFancyName/conversion_service/service/currency_manager/model"
)

const base model.CurrencyType = model.EUR

type ErrorResponse struct {
	Success bool `xml:"success,attr" json:"success"`
	Message string `xml:"message,attr" json:"message"`
}

type SucessfulResponse struct {
	Success bool `xml:"success,attr" json:"success"`
	Base string `xml:"base,attr" json:"base"`
	Input string `xml:"input,attr" json:"input"`
	Amount float64 `xml:"amount,attr" json:"amount"`
	Rates map[string]map[string]float64 `xml:"rate,attr" json:"rate"`
}

func (s *Server) handleGetCurrenciesExchangeRate(w http.ResponseWriter, r *http.Request) {	
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		rawResp, err := json.Marshal(
			&ErrorResponse{
				Success: false,
				Message: "invalid requested method",
			},
		)
		if err != nil {
			log.Printf("Failed to marshal error message. Error: %v \n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(rawResp)
		return
	}
	if r.URL.Path != apiURL {
		w.WriteHeader(http.StatusBadRequest)
		rawResp, err := json.Marshal(
			&ErrorResponse{
				Success: false,
				Message: "invalid api path",
			},
		)
		if err != nil {
			log.Printf("Failed to marshal error message. Error: %v \n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(rawResp)
		return
	}
	
	vals := r.URL.Query()
	cur := vals.Get("currency")
	amount := vals.Get("amount")
	format := vals.Get("format")

	if len(format) != 0 {
		if format == "xml" {
			w.Header().Set("Content-Type", "application/xml")
		} else if format != "json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			rawResp, err := json.Marshal(
				&ErrorResponse{
					Success: false,
					Message: "invalid requested format",
				},
			)
			if err != nil {
				log.Printf("Failed to marshal error message. Error: %v \n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(rawResp)
			return
		}
	}

	if len(amount) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		rawResp, err := json.Marshal(
			&ErrorResponse{
				Success: false,
				Message: "empty amount",
			},
		)
		if err != nil {
			log.Printf("Failed to marshal error message. Error: %v \n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(rawResp)
		return
	}
	a, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rawResp, err := json.Marshal(
			&ErrorResponse{
				Success: false,
				Message: fmt.Sprintf("invalid amount: %s", amount),
			},
		)
		if err != nil {
			log.Printf("Failed to marshal error message. Error: %v \n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(rawResp)
		return
	}

	
	if len(cur) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		rawResp, err := json.Marshal(
			&ErrorResponse{
				Success: false,
				Message: "empty input currency",
			},
		)
		if err != nil {
			log.Printf("Failed to marshal error message. Error: %v \n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(rawResp)
		return
	}
	_, present := model.AllCurrenciesTypes[model.CurrencyType(cur)]
	if ! present {
		w.WriteHeader(http.StatusBadRequest)
		rawResp, err := json.Marshal(
			&ErrorResponse{
				Success: false,
				Message: "not supported currency",
			},
		)
		if err != nil {
			log.Printf("Failed to marshal error message. Error: %v \n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(rawResp)
		return
	}

	ce, err := s.cm.GetCurrencies()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		rawResp, err := json.Marshal(
			&ErrorResponse{
				Success: false,
				Message: "failed to get currencies from the DB",
			},
		)
		if err != nil {
			log.Printf("Failed to marshal error message. Error: %v \n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(rawResp)
		return
	}

	var rawResp []byte
	if format == "xml" {
		rawResp, err = xml.Marshal(createResponseMessage(a, model.CurrencyType(cur), ce))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rawResp, err := json.Marshal(
				&ErrorResponse{
					Success: false,
					Message: "failed to get unmarshal response",
				},
			)
			if err != nil {
				log.Printf("Failed to marshal error message. Error: %v \n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(rawResp)
			return
		}
	} else {
		rawResp, err = json.Marshal(createResponseMessage(a, model.CurrencyType(cur), ce))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rawResp, err := json.Marshal(
				&ErrorResponse{
					Success: false,
					Message: "failed to get unmarshal response",
				},
			)
			if err != nil {
				log.Printf("Failed to marshal error message. Error: %v \n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(rawResp)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(rawResp)
}

func createResponseMessage(amount float64, input model.CurrencyType, rates []*model.CurrencyEntry) *SucessfulResponse {
	var i, b *model.CurrencyEntry
	for _, v := range rates {
		if i != nil && b != nil {
			break
		}
		if v.Type == input {
			i = v
		}
		if v.Type == base {
			b = v
		}
	}

	res := make(map[string]map[string]float64)
	for _, v := range rates {
		res[string(v.Type)] = map[string]float64{
			string(i.Type) : amount * i.Ratio / v.Ratio,
			string(b.Type) : amount * b.Ratio / v.Ratio,
		}

	}
	return &SucessfulResponse{
		Success: true,
		Amount: amount,
		Base: string(base),
		Input: string(input),
		Rates: res,
	}
}
