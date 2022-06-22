package web

// 登录接口验证
type (
	AuthLoginRequest struct {
		Mobile   string `form:"mobile" json:"mobile" binding:"required" label:"登录账号"`
		Password string `form:"password" json:"password" binding:"required" label:"登录密码"`
		Platform string `form:"platform" json:"platform" binding:"required,oneof=h5 ios windows mac web"`
	}

	AuthLoginResponse struct {
		Type        string `json:"type"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
)

// 注册接口验证
type (
	RegisterRequest struct {
		Nickname string `form:"nickname" json:"nickname" binding:"required,min=2,max=30" label:"账号昵称"`
		Mobile   string `form:"mobile" json:"mobile" binding:"required,len=11,phone" label:"手机号"`
		Password string `form:"password" json:"password" binding:"required,min=6,max=16" label:"密码"`
		SmsCode  string `form:"sms_code" json:"sms_code" binding:"required" label:"验证码"`
		Platform string `form:"platform" json:"platform" binding:"required,oneof=h5 ios windows mac web" label:"登录平台"`
	}

	RegisterResponse struct{}
)

// 账号找回接口验证
type (
	ForgetRequest struct {
		Mobile   string `form:"mobile" json:"mobile" binding:"required,len=11,phone" label:"手机号"`
		Password string `form:"password" json:"password" binding:"required,min=6,max=16" label:"密码"`
		SmsCode  string `form:"sms_code" json:"sms_code" binding:"required,len=6" label:"验证码"`
	}

	ForgetResponse struct{}
)

// 刷新授权接口
type (
	AuthRefreshRequest struct{}

	AuthRefreshResponse struct {
		Type        string `json:"type"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
)
