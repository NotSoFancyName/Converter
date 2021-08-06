package model

import "github.com/NotSoFancyName/conversion_service/proto"

type CurrencyType string

type CurrencyEntry struct {
	Type  CurrencyType
	Ratio float64
}

func ProtoToModel(p *proto.GetRatiosResponse) []*CurrencyEntry {
	var ce []*CurrencyEntry
	for _, v := range p.Ratios {
		ce = append(ce, &CurrencyEntry{
			Ratio: float64(v.Ratio),
			Type:  CurrencyType(v.Currency),
		})
	}
	return ce
}
