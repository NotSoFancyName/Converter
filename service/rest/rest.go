package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NotSoFancyName/conversion_service/proto"
	"github.com/NotSoFancyName/conversion_service/service/currency_manager/model"
)

const base model.CurrencyType = "EUR"

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SucessfulResponse struct {
	Success bool                          `json:"success"`
	Base    string                        `json:"base"`
	Input   string                        `json:"input"`
	Amount  float64                       `json:"amount"`
	Rates   map[string]map[string]float64 `json:"rate"`
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

	if len(format) != 0 && format != "json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		rawResp, err := json.Marshal(
			&ErrorResponse{
				Success: false,
				Message: "not supported requested format",
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
	if len(cur) != 3 {
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

	clnt := proto.NewCurrencyFetcherClient(s.cc)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := clnt.GetRatios(ctx, &proto.GetRatiosRequest{})
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
	ce := model.ProtoToModel(resp)

	rawResp, err := json.Marshal(createResponseMessage(a, model.CurrencyType(cur), ce))
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
			string(i.Type): amount * v.Ratio / i.Ratio,
			string(b.Type): amount * v.Ratio / b.Ratio,
		}

	}
	return &SucessfulResponse{
		Success: true,
		Amount:  amount,
		Base:    string(base),
		Input:   string(input),
		Rates:   res,
	}
}
