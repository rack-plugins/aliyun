package aliyun

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func help(ctx *gin.Context) {
	ctx.String(http.StatusOK, `POST /`+ID+`/addsgrule \
-H "content-type: application/json" \
-d '{
	"ip": "192.168.0.100",
	"sgid": "sg-2zeb1ux0h4683ehrocq0",
	"remark": "test",
	"policy": "[drop|accept(default)]",
}'
`)
}
