package audittrails

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/bngbngstwnd/library-go-chub/model/entity"
	"github.com/bngbngstwnd/library-go-chub/util"
	"github.com/gin-gonic/gin"
)

type AuditTrailWrapper struct {
	Data    entity.AuditTrail
	BodyLog bodyLogWriter
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func NewAuditTrails(wr gin.ResponseWriter) AuditTrailWrapper {
	at := AuditTrailWrapper{}
	at.BodyLog = bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: wr}
	return at
}

func (at *AuditTrailWrapper) GetRequestData(c *gin.Context) {
	at.Data.Time = time.Now().Format("2006-01-02T15:04:05-0700")
	at.Data.Endpoint = c.Request.URL.String()
	at.Data.Method = c.Request.Method
	at.Data.User = c.Request.Header.Get("X-CHUB-PERSONAL-NUMBER")
	at.Data.IpAddress = c.Request.Header.Get("X-Forwarded-For")
	if len(at.Data.IpAddress) == 0 {
		at.Data.IpAddress = c.Request.RemoteAddr
	}

	// inqueryKey for GET
	if at.Data.Method == "GET" {
		queryParams := c.Request.URL.Query()
		at.Data.InquiryKey = util.JSONCompactStringify(queryParams)
		at.Data.Endpoint = CleanURLFromQueryParams(at.Data.Endpoint)
	} else {
		params := c.Params
		at.Data.InquiryKey = util.JSONCompactStringify(params)
	}

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		// remove whitespaces from json
		bytesBuffer := new(bytes.Buffer)
		json.Compact(bytesBuffer, bodyBytes)
		at.Data.RequestBody = bytesBuffer.String()
	}

	// Restore the io.ReadCloser to its original state so it can be used somewhere else
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
}

func CleanURLFromQueryParams(url string) string {
	sep := "?"
	if strings.Contains(url, sep) {
		return strings.Split(url, sep)[0]
	}

	return url
}

func (at *AuditTrailWrapper) GetResponseData() {
	status := at.BodyLog.Status()
	if !(status >= 200 && status <= 299) {
		at.Data.Remark = at.BodyLog.body.String()
	} else {
		at.Data.Remark = "success"
	}
	at.Data.Status = strconv.Itoa(status)
}
