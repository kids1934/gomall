package service

import (
	"context"

	rpcuser "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"
	auth "github.com/kids1934/gomall/app/frontend/hertz_gen/frontend/auth"
	"github.com/kids1934/gomall/app/frontend/infra/rpc"
	frontendutils "github.com/kids1934/gomall/app/frontend/utils"
)

type LoginService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLoginService(Context context.Context, RequestContext *app.RequestContext) *LoginService {
	return &LoginService{RequestContext: RequestContext, Context: Context}
}

func (h *LoginService) Run(req *auth.LoginReq) (resp string, err error) {
	res, err := rpc.UserClient.Login(h.Context, &rpcuser.LoginReq{Email: req.Email, Password: req.Password})
	if err != nil {
		return
	}

	session := sessions.Default(h.RequestContext)
	session.Set("user_id", res.UserId)
	err = session.Save()
	frontendutils.MustHandleError(err)
	redirect := "/"
	if frontendutils.ValidateNext(req.Next) {
		redirect = req.Next
	}
	if err != nil {
		return "", err
	}

	return redirect, nil
}
