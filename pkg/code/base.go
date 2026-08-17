package code

import (
	"net/http"
)

var (
	OK        = Code{Code: 2000, HCode: http.StatusOK, Msg: "success"}
	Err       = Code{Code: 9999, Msg: "未定义错误"}
	ErrServer = Code{Code: 99999, Msg: "服务错误"}
	ErrField  = Code{Code: 9001, Msg: "字段错误"}
)

var (
	ErrLogin        = Code{Code: 2098, Msg: "请先登录"}
	ErrIssuedToken  = Code{Code: 2099, Msg: "令牌颁发者错误"}
	ErrInvalidToken = Code{Code: 2099, Msg: "令牌无效或过期"}

	ErrRefreshToken     = Code{Code: 2100, Msg: "令牌刷新失败"}
	ErrUpdate           = Code{Code: 2101, Msg: "更新失败"}
	ErrDelete           = Code{Code: 2102, Msg: "删除失败"}
	ErrCreate           = Code{Code: 2103, Msg: "添加失败"}
	ErrRepeat           = Code{Code: 2104, Msg: "数据已存在"}
	ErrSelect           = Code{Code: 2105, Msg: "查询失败"}
	ErrQuery            = Code{Code: 2106, Msg: "查询错误"}
	ErrRsa              = Code{Code: 2107, Msg: "请求体错误", Err: "rsa解密错误"}
	ErrNoFund           = Code{Code: 2108, Msg: "数据不存在"}
	ErrPermissionDenied = Code{Code: 2109, Msg: "权限不足"}
	ErrNil              = Code{Code: 2109, Msg: "错误"}
	ErrDecode           = Code{Code: 2110, Msg: "解析错误"}
)
