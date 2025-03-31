package global

import (
	"github.com/valyala/fasthttp"
	"main/pkg/types"
)

var AccountsList []types.AccountData
var Clients []*fasthttp.Client
var TargetProgress int64
var CurrentProgress int64 = 0
