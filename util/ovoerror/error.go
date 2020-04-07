package ovoerror

import (
	"errors"
	"fmt"
	"github.com/delanri/commonutil/util/ovostatus"
	"net/http"
	"strings"
)

type (
	ErrorStandard interface {
		Error() string
		Wrap(string)
		AppendError(string)
		GetCode() string
		GetMessage() string
		GetErrors() []string
		GetHTTPStatus() int
		IsErrorOf(string) bool
	}

	ovoError struct {
		Errors     []string
		Code       string
		Message    string
		HTTPStatus int
	}
	ResponseCode struct {
		code       string
		message    string
		httpStatus int
	}
)

const (
	SUCCESS                    = "SUCCESS"
	SYSTEM_ERROR               = "SYSTEM_ERROR"
	DUPLICATE_DATA             = "DUPLICATE_DATA"
	DATA_NOT_EXIST             = "DATA_NOT_EXIST"
	BIND_ERROR                 = "BIND_ERROR"
	RUNTIME_ERROR              = "RUNTIME_ERROR"
	DATE_NOT_VALID             = "DATE_NOT_VALID"
	VENDOR_SHUTDOWN            = "VENDOR_SHUTDOWN"
	METHOD_ARGUMENTS_NOT_VALID = "METHOD_ARGUMENTS_NOT_VALID"
	TOO_MANY_REQUEST           = "TOO_MANY_REQUEST"
	BAD_REQUEST                = "BAD_REQUEST"
	UNAUTHORIZE                = "UNAUTHORIZE"
)

var (
	responseCodes = map[string]ResponseCode{
		SUCCESS: ResponseCode{
			code:       ovostatus.SUCCESS,
			message:    "Success",
			httpStatus: http.StatusOK,
		},
		SYSTEM_ERROR: ResponseCode{
			code:       ovostatus.SYSTEM_ERROR,
			message:    "Contact our team",
			httpStatus: http.StatusInternalServerError,
		},
		DUPLICATE_DATA: ResponseCode{
			code:       ovostatus.DUPLICATE_DATA,
			message:    "Duplicate data",
			httpStatus: http.StatusOK,
		},
		DATA_NOT_EXIST: ResponseCode{
			code:       ovostatus.DATA_NOT_EXIST,
			message:    "No data exist",
			httpStatus: http.StatusOK,
		},
		BIND_ERROR: ResponseCode{
			code:       ovostatus.BIND_ERROR,
			message:    "Please fill in mandatory parameter",
			httpStatus: http.StatusOK,
		},
		RUNTIME_ERROR: ResponseCode{
			code:       ovostatus.RUNTIME_ERROR,
			message:    "Runtime Error",
			httpStatus: http.StatusInternalServerError,
		},
		DATE_NOT_VALID: ResponseCode{
			code:       ovostatus.DATE_NOT_VALID,
			message:    "Date not valid",
			httpStatus: http.StatusOK,
		},
		VENDOR_SHUTDOWN: ResponseCode{
			code:       ovostatus.VENDOR_SHUTDOWN,
			message:    "Vendor is Shutdown",
			httpStatus: http.StatusOK,
		},
		METHOD_ARGUMENTS_NOT_VALID: ResponseCode{
			code:       ovostatus.METHOD_ARGUMENTS_NOT_VALID,
			message:    "Method argument is not valid",
			httpStatus: http.StatusOK,
		},
		TOO_MANY_REQUEST: ResponseCode{
			code:       ovostatus.TOO_MANY_REQUEST,
			message:    "Invalid data",
			httpStatus: http.StatusOK,
		},
		BAD_REQUEST: ResponseCode{
			code:       ovostatus.BAD_REQUEST,
			message:    "Bad request",
			httpStatus: http.StatusBadRequest,
		},
		UNAUTHORIZE: ResponseCode{
			code:       ovostatus.UNAUTHORIZE,
			message:    "Unauthorized",
			httpStatus: http.StatusUnauthorized,
		},
	}
)

func (e ovoError) Error() string {
	err := e.Errors
	if len(err) > 0 {
		return err[0]
	} else {
		return ""
	}
}

func (e ovoError) Wrap(errMessage string) {
	e.Errors[0] = fmt.Sprintf("%s: %s", errMessage, e.Errors[0])
}

func (e *ovoError) AppendError(errMessage string) {
	e.Errors = append(e.Errors, errMessage)
}

func (e ovoError) GetCode() string {
	return e.Code
}

func (e ovoError) GetMessage() string {
	return e.Message
}

func (e ovoError) GetErrors() []string {
	return e.Errors
}

func (e ovoError) GetHTTPStatus() int {
	return e.HTTPStatus
}

func (e ovoError) IsErrorOf(code string) bool {
	if strings.ToLower(e.Code) == strings.ToLower(code) {
		return true
	}
	return false
}

func New(code string, err error) ErrorStandard {
	if code == SUCCESS {
		errCode := responseCodes[SUCCESS].code
		errMessage := responseCodes[SUCCESS].message
		errHTTPStatus := responseCodes[SUCCESS].httpStatus

		return &ovoError{
			Errors:     []string{},
			Code:       errCode,
			Message:    errMessage,
			HTTPStatus: errHTTPStatus,
		}
	}

	errCode := responseCodes[SYSTEM_ERROR].code
	errMessage := responseCodes[SYSTEM_ERROR].message
	errHTTPStatus := responseCodes[SYSTEM_ERROR].httpStatus
	errorList := make([]string, 0)

	if ovoError, ok := responseCodes[code]; ok {
		errCode = ovoError.code
		errMessage = ovoError.message
		errHTTPStatus = ovoError.httpStatus

		if err != nil {
			errorList = append(errorList, err.Error())
		}
	}

	return &ovoError{
		Errors:     errorList,
		Code:       errCode,
		Message:    errMessage,
		HTTPStatus: errHTTPStatus,
	}
}

func Wrap(err error, errMessage string) error {
	if err == nil {
		err = errors.New(errMessage)
		return err
	}

	if s, ok := err.(ErrorStandard); ok {
		s.Wrap(errMessage)
		return s
	} else {
		errTemp := errors.New(fmt.Sprintf("%s: %s", errMessage, err.Error()))
		return errTemp
	}
}

func AppendError(err error, errMessage string) error {
	if s, ok := err.(ErrorStandard); ok {
		s.AppendError(errMessage)
		return s
	}
	return err
}

func IsErrorOf(err error, code string) bool {
	if s, ok := err.(ErrorStandard); ok {
		return s.IsErrorOf(code)
	}
	return false
}
