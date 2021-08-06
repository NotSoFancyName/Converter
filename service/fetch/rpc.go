package fetch

import (
	"context"
	"log"
	"net"

	"github.com/NotSoFancyName/conversion_service/proto"
	"google.golang.org/grpc"
)

func (f *Fetcher) GetRatios(ctx context.Context, in *proto.GetRatiosRequest) (*proto.GetRatiosResponse, error) {
	ce, err := f.cm.GetCurrencies()
	if err != nil {
		return nil, err
	}

	var resp proto.GetRatiosResponse
	for _, v := range ce {
		resp.Ratios = append(resp.Ratios, &proto.GetRatiosResponse_Ratio{
			Currency: string(v.Type),
			Ratio:    float32(v.Ratio),
		})
	}
	return &resp, nil
}

func (f *Fetcher) RunRPCServer(port string, errs chan<- error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		errs <- err
	}
	s := grpc.NewServer()
	proto.RegisterCurrencyFetcherServer(s, f)

	log.Printf("gRPC server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		errs <- err
	}
}
