package code

import (
	"net/http"
)

var (
	ErrPublic = Code{Code: 999, Msg: "不可投诉或屏蔽自己"}
	ErrFilter = Code{Code: 9001, Msg: "筛选项不存在"}

	ErrOssCallback = Code{HCode: http.StatusOK, Code: 1001, Msg: "文件上传回调错误"}
	ErrFileType    = Code{Code: 2201, Msg: "文件上传类型选择错误"}
	ErrGenPolicy   = Code{Code: 2202, Msg: "文件策略签名错误"}
	ErrFileUp      = Code{Code: 2203, Msg: "文件验证错误"}

	ErrPermission = Code{Code: 2203, Msg: "请先申请敏感信息修改权限"}

	ErrCallbackPay = Code{Code: 2301, Msg: "回调验证错误"}

	ErrAccountIsNoSeller = Code{Code: 3009, Msg: "非商家账户"}
	ErrAccountNoFund     = Code{Code: 3010, Msg: "账户不存在"}
	ErrMobileCodeErr     = Code{Code: 3011, Msg: "验证码错误"}
	ErrSendMobileCodeErr = Code{Code: 3012, Msg: "验证码发送失败"}
	ErrMobileSMSRE       = Code{Code: 3013, Msg: "短信发送重试太频繁"}

	ErrNoSellerMaster             = Code{Code: 3099, Msg: "账户不是主账户"}
	ErrSellerNoInit               = Code{Code: 3100, Msg: "商户未初始化"}
	ErrSellerInit                 = Code{Code: 3101, Msg: "商户初始化失败"}
	ErrSellerInitRepeat           = Code{Code: 3102, Msg: "商户已经初始化"}
	ErrSellerNoFound              = Code{Code: 3103, Msg: "商户不存在"}
	ErrPartnerNoFound             = Code{Code: 3104, Msg: "城市合伙人不存在"}
	ErrCouponInvalid              = Code{Code: 3105, Msg: "优惠券无效"}
	ErrUserNoFound                = Code{Code: 4001, Msg: "用户不存在"}
	ErrUserBindWxChat             = Code{Code: 4001, Msg: "用户需要先绑定微信"}
	ErrUserNotIsNotMyself         = Code{Code: 4003, Msg: "用户不是本人"}
	ErrCouponAlreadyReceived      = Code{Code: 4004, Msg: "优惠券已领取"}
	ErrActivityPointsInsufficient = Code{Code: 4103, Msg: "积分不足"}

	ErrOrderNoFound = Code{Code: 5011, Msg: "订单不存在"}
	ErrOrderType    = Code{Code: 5012, Msg: "订单类型错误"}
	ErrOrderRepeat  = Code{Code: 5013, Msg: "订单已核销"}
	ErrVerifyCode   = Code{Code: 5014, Msg: "核销码错误"}
	ErrVerify       = Code{Code: 5015, Msg: "核销失败"}
	// 订单创建失败
	ErrOrderCreate           = Code{Code: 5016, Msg: "订单创建失败"}
	ErrOrderNoFoundOrTimeout = Code{Code: 5017, Msg: "订单不存在或超时"}
	ErrPayWay                = Code{Code: 5018, Msg: "支付方式选择错误"}
	ErrPayWayNotAvailable    = Code{Code: 5019, Msg: "支付方式不可用"}
	ErrPayOrderGen           = Code{Code: 5020, Msg: "支付订单生成失败"}
	ErrOrderCanceledFailed   = Code{Code: 5021, Msg: "订单无法取消状态错误"}
	ErrOrderNotCompleted     = Code{Code: 5022, Msg: "订单未完成"}
	ErrOrderCommented        = Code{Code: 5023, Msg: "订单已评论"}
	ErrOrderCommentNoFound   = Code{Code: 5024, Msg: "订单未评论"}
	ErrOrderDeliveryMethod   = Code{Code: 5025, Msg: "请选择有效的收货渠道"}
	ErrOrderAddressEmpty     = Code{Code: 5026, Msg: "地址不能为空"}
	ErrOrderAddressNoFound   = Code{Code: 5027, Msg: "地址不存在"}

	// 该订单售后已提交
	ErrOrderRefundRepeat = Code{Code: 5100, Msg: "该订单售后已提交"}
	// 该订单无法申请售后
	ErrOrderRefundApply = Code{Code: 5101, Msg: "该订单无法申请售后"}
	// 订单售后申请失败
	ErrOrderRefundApplyFailed  = Code{Code: 5102, Msg: "订单售后申请失败"}
	ErrOrderRefundApplyNoFound = Code{Code: 5103, Msg: "订单售后申请不存在"}
	ErrOrderRefundStatusCancel = Code{Code: 5104, Msg: "当前状态不可取消"}

	ErrResClass = Code{Code: 5101, Msg: "资源类目无效"}

	ErrLike   = Code{Code: 7101, Msg: "点赞重复"}
	ErrUnLike = Code{Code: 7102, Msg: "点赞已取消"}

	// 商品不存在
	ErrProductNoFound    = Code{Code: 8001, Msg: "商品不存在"}
	ErrProductSkuNoFound = Code{Code: 8103, Msg: "商品规格不存在"}
	ErrActivityNoFound   = Code{Code: 8002, Msg: "活动不存在"}
	ErrProductNoShow     = Code{Code: 8104, Msg: "该商品未上架"}
	ErrProductSkuStock   = Code{Code: 8105, Msg: "商品规格无库存"}

	ErrActNoFund = Code{Code: 8002, Msg: "活动不存在"}
	ErrActJoin   = Code{Code: 8003, Msg: "活动参加失败"}
)
