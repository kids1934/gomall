package service

import (
	"context"

	cart "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/kids1934/gomall/app/cart/biz/dal/mysql"
	"github.com/kids1934/gomall/app/cart/biz/model"
)

type EmptyCartService struct {
	ctx context.Context
} // NewEmptyCartService new EmptyCartService
func NewEmptyCartService(ctx context.Context) *EmptyCartService {
	return &EmptyCartService{ctx: ctx}
}

// Run create note info
func (s *EmptyCartService) Run(req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	// Finish your business logic.
	err = model.EmptyCart(mysql.DB, s.ctx, req.GetUserId())
	if err != nil {
		return &cart.EmptyCartResp{}, kerrors.NewBizStatusError(50001, "empty cart error")
	}

	return &cart.EmptyCartResp{}, nil
}
