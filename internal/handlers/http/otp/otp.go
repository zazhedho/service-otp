package handlerotp

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"service-sender/internal/dto"
	interfaceotp "service-sender/internal/interfaces/otp"
	serviceotp "service-sender/internal/services/otp"
	"service-sender/pkg/logger"
	"service-sender/pkg/messages"
	"service-sender/pkg/response"
	"service-sender/utils"
)

type HandlerOTP struct {
	Service interfaceotp.ServiceOTPInterface
}

func NewOTPHandler(s interfaceotp.ServiceOTPInterface) *HandlerOTP {
	return &HandlerOTP{Service: s}
}

func (h *HandlerOTP) SendRegisterOTP(ctx *gin.Context) {
	var req dto.OTPSendRequest
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[OTPHandler][SendRegisterOTP]"

	if err := ctx.BindJSON(&req); err != nil {
		logger.WriteLogWithContext(ctx, logger.LogLevelError, logPrefix+"; BindJSON ERROR: "+err.Error())
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	appName := strings.TrimSpace(ctx.GetHeader("X-App-Name"))
	if appName == "" {
		appName = strings.TrimSpace(utils.GetEnv("OTP_APP_NAME", "Account Verification").(string))
	}
	logger.WriteLogWithContext(ctx, logger.LogLevelDebug, fmt.Sprintf("%s; Request email: %s; AppName: %s;", logPrefix, req.Email, appName))

	err := h.Service.SendRegisterOTP(ctx.Request.Context(), req.Email, appName)
	if err != nil {
		logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; Service.SendRegisterOTP error: %v", logPrefix, err))
		if throttle := new(serviceotp.ThrottleError); errors.As(err, &throttle) {
			retryAfter := int(throttle.RetryAfter.Seconds())
			if retryAfter > 0 {
				ctx.Header("Retry-After", strconv.Itoa(retryAfter))
			}
			res := response.Response(http.StatusTooManyRequests, messages.MsgFail, logId, nil)
			res.Error = response.Errors{Code: http.StatusTooManyRequests, Message: "Please wait before requesting another OTP"}
			ctx.JSON(http.StatusTooManyRequests, res)
			return
		}

		if errors.Is(err, serviceotp.ErrOTPNotConfigured) || errors.Is(err, serviceotp.ErrOTPDeliveryFailed) {
			res := response.Response(http.StatusInternalServerError, messages.MsgFail, logId, nil)
			res.Error = response.Errors{Code: http.StatusInternalServerError, Message: "OTP service is not available"}
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}

		res := response.Response(http.StatusBadRequest, messages.MsgFail, logId, nil)
		res.Error = response.Errors{Code: http.StatusBadRequest, Message: "Unable to process OTP request"}
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response(http.StatusOK, messages.MsgSuccess, logId, map[string]string{"email": req.Email})
	logger.WriteLogWithContext(ctx, logger.LogLevelInfo, fmt.Sprintf("%s; OTP sent to: %s", logPrefix, req.Email))
	ctx.JSON(http.StatusOK, res)
}

func (h *HandlerOTP) VerifyRegisterOTP(ctx *gin.Context) {
	var req dto.OTPVerifyRequest
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[OTPHandler][VerifyRegisterOTP]"

	if err := ctx.BindJSON(&req); err != nil {
		logger.WriteLogWithContext(ctx, logger.LogLevelError, logPrefix+"; BindJSON ERROR: "+err.Error())
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	logger.WriteLogWithContext(ctx, logger.LogLevelDebug, fmt.Sprintf("%s; Request email: %s;", logPrefix, req.Email))

	err := h.Service.VerifyRegisterOTP(ctx.Request.Context(), req.Email, req.Code)
	if err != nil {
		logger.WriteLogWithContext(ctx, logger.LogLevelWarn, fmt.Sprintf("%s; Service.VerifyRegisterOTP error: %v", logPrefix, err))
		res := response.Response(http.StatusBadRequest, messages.MsgFail, logId, nil)
		res.Error = response.Errors{Code: http.StatusBadRequest, Message: "OTP verification failed"}
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response(http.StatusOK, messages.MsgSuccess, logId, map[string]string{"email": req.Email})
	logger.WriteLogWithContext(ctx, logger.LogLevelInfo, fmt.Sprintf("%s; OTP verified for: %s", logPrefix, req.Email))
	ctx.JSON(http.StatusOK, res)
}
