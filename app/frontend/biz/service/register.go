package service

import (
	"context"

	rpcuser "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"
	auth "github.com/kids1934/gomall/app/frontend/hertz_gen/frontend/auth"
	common "github.com/kids1934/gomall/app/frontend/hertz_gen/frontend/common"
	"github.com/kids1934/gomall/app/frontend/infra/rpc"
)

type RegisterService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRegisterService(Context context.Context, RequestContext *app.RequestContext) *RegisterService {
	return &RegisterService{RequestContext: RequestContext, Context: Context}
}

func (h *RegisterService) Run(req *auth.RegisterReq) (resp *common.Empty, err error) {
	res, err := rpc.UserClient.Register(h.Context, &rpcuser.RegisterReq{
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.Password,
	})
	if err != nil {
		return nil, err
	}

	session := sessions.Default(h.RequestContext)
	session.Set("user_id", res.UserId)
	err = session.Save()
	if err != nil {
		return nil, err
	}
	return
}
