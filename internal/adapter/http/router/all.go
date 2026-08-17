package router

import (
	gc "github.com/1303-yzym/MoonshotWell/internal/adapter/http/context"
	"github.com/1303-yzym/MoonshotWell/internal/application/service"
	"github.com/1303-yzym/MoonshotWell/internal/infrastructure/state"
	"github.com/1303-yzym/MoonshotWell/pkg/warp"
	"github.com/gin-gonic/gin"
)

type InitRouter[C warp.Context[*gin.Context]] interface {
	InitRouter(gr warp.RouterGroup[*gin.Context, *gc.Context])
}

func registerRoutingGroup(gr warp.RouterGroup[*gin.Context, *gc.Context], rs ...InitRouter[warp.Context[*gin.Context]]) {
	for _, r := range rs {
		r.InitRouter(gr)
	}
}

func SetupRouter(gr warp.RouterGroup[*gin.Context, *gc.Context], appState *state.AppState, appService *service.Service) {
	// 在这创建全局的处理器
	// hdr := handler.New(appState, appService)

	registerRoutingGroup(gr) //Base{hdr},
	//// 用户
	//User{hdr},
	//UserAddress{hdr},
	//UserBlock{hdr},
	//UserGroup{hdr},
	//UserFollow{hdr},
	//UserAccount{hdr},
	//UserNotification{hdr},
	//UserMember{hdr},
	//PushMessage{hdr},
	//Points{hdr},
	//// 活动
	//Circle{hdr},
	//ActivityFeed{hdr},
	//Activity{hdr},
	//Feed{hdr},
	//ActivityLottery{hdr},
	//// 商城
	//Mall{hdr},
	//MallProduct{hdr},
	//MallOrder{hdr},
	//Resource{hdr},
	//// 上传
	//Upload{hdr},
	//// 回调
	//Callback{hdr},
	//// 奖品
	//Prize{hdr},
	//// 商户
	//Seller{hdr},
	//// 福利
	//Welfare{hdr},
	//// 主动事件
	//Event{hdr},
	//// 客服
	//KF{hdr},
	//Admin{hdr},
	//// 优惠券
	//Coupon{hdr},
	//Other{hdr},

}
