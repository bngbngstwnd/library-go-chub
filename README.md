# library-go-chub

CHUB Library

CHUB Library is private library, consists of standardized functions and constants which is used by Customer Hub Project.
Installation

    Set GOPRIVATE variable in Go env to allow using private repository;

    go env -w GOPRIVATE=github.com/bngbngstwnd/library-go-chub

    Set git configuration, replacing "https://github.com" into "https://GITHUB_KEY@github.com". Ask administrator for GITHUB_KEY value;

    git config --global url."https://*GITHUB_KEY*@github.com".insteadOf "https://github.com"

    Get the library.

    go get -u github.com/bngbngstwnd/library-go-chub

#Usage

A. Constant
1. Response Code, Response Message and Error Code

These collections of constants are being used as standard response for Customer Hub
Response Code Library

RESPONSE_MESSAGE_SUCCESS                        = "Success"
RESPONSE_MESSAGE_INVALID_BODY_REQ               = "Invalid Body Request"
RESPONSE_MESSAGE_INVALID_HEADER_REQ             = "Invalid Header Request"
RESPONSE_MESSAGE_INVALID_QUERY_PARAMS           = "Invalid Query Parameters"
RESPONSE_MESSAGE_BODY_REQ_EMPTY                 = "Body Request Empty"
RESPONSE_MESSAGE_INVALID_DATA_TYPE              = "Invalid Data Type"
RESPONSE_MESSAGE_INVALID_DATE_FORMAT            = "Invalid Date Format"
RESPONSE_MESSAGE_INVALID_TIME_FORMAT            = "Invalid Time Format"
RESPONSE_MESSAGE_INVALID_AUTH_TOKEN             = "Invalid Auth Token"
RESPONSE_MESSAGE_INVALID_DATA_ALREADY_EXIST     = "Data Already Exist"
RESPONSE_MESSAGE_AUTH_TOKEN_EXPIRED             = "Auth Token Expired"
RESPONSE_MESSAGE_AUTH_TOKEN_EMPTY               = "Auth Token Empty"
RESPONSE_MESSAGE_DATA_NOT_FOUND                 = "Data Not Found"
RESPONSE_MESSAGE_ROUTE_NOT_FOUND                = "Route Not Found"
RESPONSE_MESSAGE_DATABASE_ERROR                 = "Database Error"
RESPONSE_MESSAGE_PRODUCE_MESSAGE_ERROR          = "Failed to produce message"
RESPONSE_MESSAGE_TIMEOUT                        = "Timeout"
RESPONSE_MESSAGE_UNDEFINED_ERROR                = "Undefined Error"
RESPONSE_MESSAGE_GENERAL_ERROR                  = "General Error"

Response Message Library

RESPONSE_CODE_SUCCESS           = "00"
RESPONSE_CODE_BAD_REQUEST       = "01"
RESPONSE_CODE_AUTH_ERROR        = "02"
RESPONSE_CODE_NOT_FOUND         = "03"
RESPONSE_CODE_INTERNAL_ERROR    = "05"

Error Code Library

ERROR_CODE_SUCCESS                      = "000"
ERROR_CODE_INVALID_BODY_REQUEST         = "010"
ERROR_CODE_BODY_REQUEST_EMPTY           = "011"
ERROR_CODE_INVALID_DATA_TYPE            = "012"
ERROR_CODE_INVALID_DATE_FORMAT          = "013"
ERROR_CODE_INVALID_TIME_FORMAT          = "014"
ERROR_CODE_INVALID_HEADER_REQUEST       = "015"
ERROR_CODE_INVALID_QUERY_PARAMS         = "016"
ERROR_CODE_INVALID_DATA_ALREADY_EXIST   = "017"
ERROR_CODE_INVALID_AUTH_TOKEN           = "020"
ERROR_CODE_AUTH_TOKEN_EXPIRED           = "021"
ERROR_CODE_AUTH_TOKEN_EMPTY             = "022"
ERROR_CODE_DATA_NOT_FOUND               = "030"
ERROR_CODE_DATABASE_ERROR               = "050"
ERROR_CODE_TIMEOUT                      = "059"
ERROR_CODE_PRODUCE_MESSAGE_ERROR        = "060"
ERROR_CODE_UNDEFINED_ERROR              = "888"
ERROR_CODE_GENERAL_ERROR                = "999"

Examples

import (
    "github.com/bngbngstwnd/library-go-chub/constant"
)

// returns 'Success'
constant.RESPONSE_MESSAGE_SUCCESS

// returns '00'
constant.RESPONSE_CODE_SUCCESS

// returns '000'
constant.ERROR_CODE_SUCCESS

2. Date and Time Format

Constants that contains Date and Time Format which can be used when parsing to Date and Time data type
Date Layout Library

    DATE_LAYOUT             = "2006-01-02"
    DATE_LAYOUT_DDMMYYYY    = "02-01-2006"

Date Time Layout Library

    DATE_TIME_LAYOUT            = "2006-01-02T15:04:05Z"
    DATE_TIME_LAYOUT_NON_TZ     = "2006-01-02 15:04:05"

Example

    import(
        "bitbucket.org/bridce/ms-chub-prodrec/utils"
    )

    layoutedTime, err := time.Parse(utils.DATE_TIME_LAYOUT_NON_TZ, "2021-10-19 10:30:45")

B. Request and Response Model

Models Collection that Customer Hub use as Request and Response Structures
Telegram Request Model

    type TelegramNotificationRequest struct {
        ChannelID string                      `json:"channel_id"`
        Payload   TelegramNotificationPayload `json:"payload"`
    }

    // TelegramNotificationPayload : Parameter payload notifikasi ke Telegram
    type TelegramNotificationPayload struct {
        Message string `json:"message"`
        Apps    string `json:"apps"`
        Status  string `json:"status"`
    }

Example


Response Model

    type ResponseContainer struct {
        Response Response `json:"RESPONSE"`
    }

    type Response struct {
        StatusCode      int           `json:"STATUS_CODE"`
        ErrorCode       *string       `json:"ERROR_CODE"`
        ResponseCode    *string       `json:"RESPONSE_CODE"`
        ResponseMessage *string       `json:"RESPONSE_MESSAGE"`
        Errors          []string      `json:"ERRORS"`
        Data            interface{}   `json:"DATA"`
        Info            *ResponseInfo `json:"INFO,omitempty"`
    }

    type ErrorContainer struct {
        Response Response `json:"RESPONSE"`
    }

    type ResponseInfo struct {
        Limit    int `json:"LIMIT"`
        Page     int `json:"PAGE"`
        PageSize int `json:"PAGE_SIZE"`
        Total    int `json:"TOTAL"`
    }

Example

    import(
        "github.com/bngbngstwnd/library-go-chub/model/response"
    )

    info := response.ResponseInfo{
        Limit:    limit,
        Page:     page,
        PageSize: pageSize,
        Total:    total,
    }

Response Builder Model

    func BuildSuccessResponse(data interface{}) *ResponseContainer {
        return &ResponseContainer{
            Response: Response{
                StatusCode:      http.StatusOK,
                ErrorCode:       &constant.ERROR_CODE_SUCCESS,
                ResponseCode:    &constant.RESPONSE_CODE_SUCCESS,
                ResponseMessage: &constant.RESPONSE_MESSAGE_SUCCESS,
                Errors:          nil,
                Data:            data,
                Info:            nil,
            },
        }
    }

    func BuildSuccessResponseWithInfo(data interface{}, info *ResponseInfo) *ResponseContainer {
        return &ResponseContainer{
            Response: Response{
                StatusCode:      http.StatusOK,
                ErrorCode:       &constant.ERROR_CODE_SUCCESS,
                ResponseCode:    &constant.RESPONSE_CODE_SUCCESS,
                ResponseMessage: &constant.RESPONSE_MESSAGE_SUCCESS,
                Errors:          nil,
                Data:            data,
                Info:            info,
            },
        }
    }

    func BuildDataNotFoundResponse() *ErrorContainer {
        return &ErrorContainer{
            Response: Response{
                StatusCode:      http.StatusNotFound,
                ErrorCode:       &constant.ERROR_CODE_DATA_NOT_FOUND,
                ResponseCode:    &constant.RESPONSE_CODE_NOT_FOUND,
                ResponseMessage: &constant.RESPONSE_MESSAGE_DATA_NOT_FOUND,
                Errors:          nil,
                Data:            nil,
                Info:            nil,
            },
        }
    }

    func BuildDataNotFoundResponseWithMessage(msg string) *ErrorContainer {
        return &ErrorContainer{
            Response: Response{
                StatusCode:      http.StatusNotFound,
                ErrorCode:       &constant.ERROR_CODE_DATA_NOT_FOUND,
                ResponseCode:    &constant.RESPONSE_CODE_NOT_FOUND,
                ResponseMessage: &constant.RESPONSE_MESSAGE_DATA_NOT_FOUND,
                Errors:          strings.Split(msg, "\n"),
                Data:            nil,
                Info:            nil,
            },
        }
    }

    func BuildBadRequestResponse(errCode, respCode, errMessage, throwable string) *ErrorContainer {
        return &ErrorContainer{
            Response: Response{
                StatusCode:      http.StatusBadRequest,
                ErrorCode:       &errCode,
                ResponseCode:    &respCode,
                ResponseMessage: &errMessage,
                Errors:          strings.Split(throwable, "\n"),
                Data:            nil,
                Info:            nil,
            },
        }
    }

    func BuildInternalErrorResponse(errCode, respCode, errMessage, throwable string) *ErrorContainer {
        return &ErrorContainer{
            Response: Response{
                StatusCode:      http.StatusInternalServerError,
                ErrorCode:       &errCode,
                ResponseCode:    &respCode,
                ResponseMessage: &errMessage,
                Errors:          strings.Split(throwable, "\n"),
                Data:            nil,
                Info:            nil,
            },
        }
    }

    func BuildRouteNotFoundResponse() *ErrorContainer {
        return &ErrorContainer{
            Response: Response{
                StatusCode:      http.StatusNotFound,
                ErrorCode:       &constant.ERROR_CODE_DATA_NOT_FOUND,
                ResponseCode:    &constant.RESPONSE_CODE_NOT_FOUND,
                ResponseMessage: &constant.RESPONSE_MESSAGE_ROUTE_NOT_FOUND,
                Errors:          nil,
                Data:            nil,
                Info:            nil,
            },
        }
    }

    func BuildEmptyBodyReqResponse(errMessage, throwable string) *ErrorContainer {
        return &ErrorContainer{
            Response: Response{
                StatusCode:      http.StatusBadRequest,
                ErrorCode:       &constant.ERROR_CODE_BODY_REQUEST_EMPTY,
                ResponseCode:    &constant.RESPONSE_CODE_BAD_REQUEST,
                ResponseMessage: &errMessage,
                Errors:          strings.Split(throwable, "\n"),
                Data:            nil,
            },
        }
    }

    func BuildInvalidTypeResponse(errMessage, throwable string) *ErrorContainer {
        return &ErrorContainer{
            Response: Response{
                StatusCode:      http.StatusBadRequest,
                ErrorCode:       &constant.ERROR_CODE_INVALID_DATA_TYPE,
                ResponseCode:    &constant.RESPONSE_CODE_BAD_REQUEST,
                ResponseMessage: &errMessage,
                Errors:          strings.Split(throwable, "\n"),
                Data:            nil,
            },
        }
    }

Example

    import(
        "github.com/bngbngstwnd/library-go-chub/constant"
        "github.com/bngbngstwnd/library-go-chub/model/response"
    )

    c *gin.Context

    errResp := response.BuildBadRequestResponse(constant.ERROR_CODE_INVALID_HEADER_REQUEST, constant.RESPONSE_CODE_BAD_REQUEST, constant.RESPONSE_MESSAGE_INVALID_HEADER_REQ, errPn.Error())

    c.JSON(http.StatusBadRequest, errResp)

C. Telegram Notifications

Package which can be used to send message throguh telegram

    type Notification interface {
        Send(status, message string) error
    }

    type telegram struct {
        host      string
        channelID string
        appName   string
    }

    func TelegramNotification(host, channelID, appName string) Notification {
        return &telegram{
            host:      host,
            channelID: channelID,
            appName:   appName,
        }
    }

    // Send : Digunakan untuk mengirim notifikasi ke telegram
    func (tele *telegram) Send(status, message string) error {

        url := tele.host + "/broadcast"
        channelID := tele.channelID
        applicationName := tele.appName

        var payload request.TelegramNotificationPayload
        payload.Apps = applicationName
        payload.Message = message
        payload.Status = status

        var param request.TelegramNotificationRequest
        param.ChannelID = channelID
        param.Payload = payload
        s, _ := json.Marshal(param)

        res, err := http.Post(url, "application/json", bytes.NewBuffer(s))

        if err != nil {
            log.Println("Failed sending notification: ", err.Error())
            return err
        }
        defer res.Body.Close()

        log.Println("Telegram status:", res.Status)

        return nil
    }

Example


D. Utils
1. Array Functions

    func IsInArrayInt(arr []int, el int) bool {
        for _, e := range arr {
            if e == el {
                return true
            }
        }
        return false
    }

    func IsInArrayStr(arr []string, el string) bool {
        for _, e := range arr {
            if e == el {
                return true
            }
        }
        return false
    }

    func IsIntersectArrayStr(arr1, arr2 []string) bool {
        for _, e := range arr1 {
            if IsInArrayStr(arr2, e) {
                return true
            }
        }
        return false
    }

    func IsIntersectArrayInt(arr1, arr2 []int) bool {
        for _, e := range arr1 {
            if IsInArrayInt(arr2, e) {
                return true
            }
        }
        return false
    }

    // this function returns non-intersect element from first array
    func FilterOutIntersectSliceStr(arr1, arr2 []string) []string {
        var result []string
        for _, e := range arr1 {
            if IsInArrayStr(arr2, e) {
                continue
            }
            result = append(result, e)
        }
        return result
    }

2. JSON Function

    func JSONCompactStringify(data interface{}) string {
        bytesOfObj, _ := json.Marshal(data)
        bytesBuffer := new(bytes.Buffer)
        json.Compact(bytesBuffer, bytesOfObj)
        return bytesBuffer.String()
    }

3. Time Function

    func GetCurrentSecond(input time.Time) int {
        return (input.Hour() * 3600) + input.Minute()*60 + input.Second()
    }

Example


E. Validator

    var (
        v     *validator.Validate
        trans ut.Translator
    )

    func GetValidator() (*validator.Validate, ut.Translator) {
        return v, trans
    }

    func InitValidator() {
        translator := en.New()
        uni := ut.New(translator, translator)

        // this is usually known or extracted from http 'Accept-Language' header
        var found bool
        trans, found = uni.GetTranslator("en")
        if !found {
            fmt.Println("translator not found")
        }

        v = validator.New()

        if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
            fmt.Println(err)
        }

        _ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
            return ut.Add("required", "{0} tidak boleh kosong", true)
        }, func(ut ut.Translator, fe validator.FieldError) string {
            t, _ := ut.T("required", fe.Field())
            return t
        })

        _ = v.RegisterTranslation("numeric", trans, func(ut ut.Translator) error {
            return ut.Add("numeric", "{0} harus berupa angka", true)
        }, func(ut ut.Translator, fe validator.FieldError) string {
            t, _ := ut.T("numeric", fe.Field())
            return t
        })

        v.RegisterTagNameFunc(func(fld reflect.StructField) string {
            name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
            if name == "-" {
                return ""
            }
            return name
        })

        _ = v.RegisterValidation("campaigndate", func(fl validator.FieldLevel) bool {
            val := fl.Field().String()
            _, err := time.Parse(constant.DATE_TIME_LAYOUT_NON_TZ, val)
            return err == nil
        })

        _ = v.RegisterTranslation("campaigndate", trans, func(ut ut.Translator) error {
            return ut.Add("campaigndate", "{0} must be in format 'yyyy-mm-dd hh:mm:ss'", true)
        }, func(ut ut.Translator, fe validator.FieldError) string {
            t, _ := ut.T("campaigndate", fe.Field())
            return t
        })

        _ = v.RegisterValidation("datedob", func(fl validator.FieldLevel) bool {
            val := fl.Field().String()
            _, err := time.Parse(constant.DATE_LAYOUT, val)
            return err == nil
        })

        _ = v.RegisterTranslation("datedob", trans, func(ut ut.Translator) error {
            return ut.Add("datedob", "{0} must be in format 'yyyy-mm-dd'", true)
        }, func(ut ut.Translator, fe validator.FieldError) string {
            t, _ := ut.T("datedob", fe.Field())
            return t
        })
    }

    func GetAndValidatePersonalNumber(c *gin.Context) (string, error) {
        personalNumber := c.GetHeader("X-CHUB-PERSONAL-NUMBER")
        if len(personalNumber) != 8 {
            return "", errors.New("invalid personal number")
        }
        if _, err := strconv.ParseInt(personalNumber, 10, 64); err != nil {
            return "", errors.New("personal number should be numeric")
        }
        return personalNumber, nil
    }

    func Validate(req interface{}) error {
        v, trans := GetValidator()
        var errorMessages []string

        errData := v.Struct(req)
        if errData != nil {
            for _, e := range errData.(validator.ValidationErrors) {
                errorMessages = append(errorMessages, e.Translate(trans))
            }
            errMessage := strings.Join(errorMessages, "\n")

            return errors.New(errMessage)
        }

        return nil
    }

Example

    import(
        "github.com/bngbngstwnd/library-go-chub/validator"
    )

    c *gin.Context

    personalNumber, errPn := validator.GetAndValidatePersonalNumber(c)
    if errPn != nil {
        errResp := response.BuildBadRequestResponse(constant.ERROR_CODE_INVALID_HEADER_REQUEST, constant.RESPONSE_CODE_BAD_REQUEST, constant.RESPONSE_MESSAGE_INVALID_HEADER_REQ, errPn.Error())
        c.JSON(http.StatusBadRequest, errResp)
        return
    }

Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
