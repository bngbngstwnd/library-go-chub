package timeout

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bngbngstwnd/library-go-chub/constant"
	"github.com/bngbngstwnd/library-go-chub/model/response"
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm"
)

func Middleware(timeMilis int) gin.HandlerFunc {

	log.Printf("initialized maximum processing time: %d milliseconds\n", timeMilis)

	return TimeoutHandler(
		time.Duration(timeMilis)*time.Millisecond,
		http.StatusRequestTimeout,
		response.BuildTimeoutResponse("Processing time is too long"))
}

type DeliveryProcess func(*gin.Context) (int, interface{}, bool)

func DeliveryWrapper(c *gin.Context, apmTx *apm.Transaction, process DeliveryProcess) {

	defer apmTx.End()
	waktuMulaiProcessing := time.Now()
	defer func() {
		lagTime := (int(time.Now().UnixNano() - waktuMulaiProcessing.UnixNano())) / 1000000
		apmTx.Context.SetLabel("lag", lagTime)
	}()

	var deliveryResponseCode int
	var deliveryOutputData interface{}
	var isSuccess bool

	ctx := c.Request.Context()

	doneChan := make(chan bool)

	go func() {
		deliveryResponseCode, deliveryOutputData, isSuccess = process(c)
		close(doneChan)
	}()

	select {

	case <-ctx.Done():
		apmTx.Result = strconv.Itoa(http.StatusRequestTimeout)
		apmTx.Context.SetLabel("status_code", http.StatusRequestTimeout)
		apmTx.Context.SetLabel("error_code", constant.ERROR_CODE_TIMEOUT)
		apmTx.Context.SetLabel("response_code", constant.RESPONSE_CODE_INTERNAL_ERROR)
		e := apm.DefaultTracer.NewError(errors.New(constant.RESPONSE_MESSAGE_TIMEOUT))
		e.SetTransaction(apmTx)
		e.Send()

		return

	case <-doneChan:
		if !isSuccess {
			if errResp, ok := deliveryOutputData.(*response.ErrorContainer); ok {
				apmTx.Result = strconv.Itoa(errResp.Response.StatusCode)
				apmTx.Context.SetLabel("status_code", errResp.Response.StatusCode)
				apmTx.Context.SetLabel("error_code", *errResp.Response.ErrorCode)
				apmTx.Context.SetLabel("response_code", *errResp.Response.ResponseCode)
				e := apm.DefaultTracer.NewError(errors.New(*errResp.Response.ResponseMessage))
				e.SetTransaction(apmTx)
				e.Send()
			}
		} else {
			if resp, ok := deliveryOutputData.(*response.ResponseContainer); ok {
				apmTx.Result = strconv.Itoa(resp.Response.StatusCode)
				apmTx.Context.SetLabel("status_code", resp.Response.StatusCode)
				apmTx.Context.SetLabel("error_code", *resp.Response.ErrorCode)
				apmTx.Context.SetLabel("response_code", *resp.Response.ResponseCode)
			}
		}

		c.JSON(deliveryResponseCode, deliveryOutputData)
	}
}

func TimeoutHandler(timeout time.Duration, responseCodeTimeout int, responseBodyTimeout interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(responseCodeTimeout, responseBodyTimeout)
			c.Abort()
		}

	}
}
