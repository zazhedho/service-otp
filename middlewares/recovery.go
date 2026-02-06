package middlewares

import (
	"fmt"
	"net/http"
	"service-sender/pkg/logger"
	"service-sender/pkg/messages"
	"service-sender/pkg/response"
	"service-sender/utils"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context, err any) {
	logId := utils.GenerateLogId(c)
	logger.WriteLogWithContext(c, logger.LogLevelPanic, fmt.Sprintf("RECOVERY; Error: %+v;", err))

	res := response.Response(http.StatusInternalServerError, fmt.Sprintf("%s (%s)", messages.MsgFail, logId.String()), logId, nil)
	c.AbortWithStatusJSON(http.StatusInternalServerError, res)
	return
}
