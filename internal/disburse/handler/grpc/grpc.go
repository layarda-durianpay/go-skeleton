package grpchandler

import (
	"context"

	"github.com/layarda-durianpay/go-skeleton/internal/disburse/app"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/app/command"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/grpcerr"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/protogen"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	app *app.Application
}

func NewGrpcServer(
	application *app.Application,
) GRPCServer {

	return GRPCServer{app: application}
}

func (g GRPCServer) Disburse(ctx context.Context, req *protogen.DisburseRequest) (*emptypb.Empty, error) {
	err := g.app.Commands.Disburse.Handle(ctx, &command.DisburseParam{
		Amount: req.Amount,
	})
	if err != nil {
		return nil, grpcerr.TransformToGRPCErr(err)
	}

	return &emptypb.Empty{}, nil
}
