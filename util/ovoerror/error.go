package ovoerror

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/delanri/commonutil/util/ovostatus"
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
	SUCCESS                     = "SUCCESS"
	SYSTEM_ERROR                = "SYSTEM_ERROR"
	DUPLICATE_DATA              = "DUPLICATE_DATA"
	DATA_NOT_EXIST              = "DATA_NOT_EXIST"
	BIND_ERROR                  = "BIND_ERROR"
	RUNTIME_ERROR               = "RUNTIME_ERROR"
	DATE_NOT_VALID              = "DATE_NOT_VALID"
	VENDOR_SHUTDOWN             = "VENDOR_SHUTDOWN"
	METHOD_ARGUMENTS_NOT_VALID  = "METHOD_ARGUMENTS_NOT_VALID"
	TOO_MANY_REQUEST            = "TOO_MANY_REQUEST"
	BAD_REQUEST                 = "BAD_REQUEST"
	UNAUTHORIZE                 = "UNAUTHORIZE"
	DUPLICATE_BANK_ACCOUNT      = "DUPLICATE_BANK_ACCOUNT"
	DIFFERENT_NAME_BANK_ACCOUNT = "DIFFERENT_NAME_BANK_ACCOUNT"
	TOP_UP_OVERLIMIT            = "TOP_UP_OVERLIMIT"
	BANK_ACCOUNT_NOT_FOUND      = "BANK_ACCOUNT_NOT_FOUND"
	ORS_BLOCKED                 = "ORS_BLOCKED"
	INVALID_DATA                = "INVALID_DATA"
	THIRD_PARTY_ERROR           = "THIRD_PARTY_ERROR"
	INELIGIBLE                  = "INELIGIBLE"
	PAYMENT_FAILED              = "PAYMENT_FAILED"
	PAYMENT_TIMED_OUT           = "PAYMENT_TIMED_OUT"
	INSUFFICIENT_FUND           = "INSUFFICIENT_FUND"
	INVALID_TRANSACTION         = "INVALID_TRANSACTION"
	ORDER_CREATION_FAILED       = "ORDER_CREATION_FAILED"
	INCOMPLETE_ACCOUNT          = "INCOMPLETE_ACCOUNT"
	PROMO_CODE_ERROR            = "PROMO_CODE_ERROR"
	PROMO_BOOK_ERROR            = "PROMO_BOOK_ERROR"
	INSUFFICIENT_APP_VERSION    = "INSUFFICIENT_APP_VERSION"
	PENDING_TRANSACTION_ERROR   = "PENDING_TRANSACTION_ERROR"
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
		DUPLICATE_BANK_ACCOUNT: ResponseCode{
			code:       ovostatus.DUPLICATE_BANK_ACCOUNT,
			message:    "Bank account is already registered",
			httpStatus: http.StatusBadRequest,
		},
		DIFFERENT_NAME_BANK_ACCOUNT: ResponseCode{
			code:       ovostatus.DIFFERENT_NAME_BANK_ACCOUNT,
			message:    "Different name with bank account name",
			httpStatus: http.StatusBadRequest,
		},
		TOP_UP_OVERLIMIT: ResponseCode{
			code:       ovostatus.TOP_UP_OVERLIMIT,
			message:    "User OVO Cash balance exceeds limit",
			httpStatus: http.StatusBadRequest,
		},
		BANK_ACCOUNT_NOT_FOUND: ResponseCode{
			code:       ovostatus.BANK_ACCOUNT_NOT_FOUND,
			message:    "Bank Account Not Found",
			httpStatus: http.StatusBadRequest,
		},
		ORS_BLOCKED: ResponseCode{
			code:       ovostatus.ORS_BLOCKED,
			message:    "Bank Account Not Found",
			httpStatus: http.StatusBadRequest,
		},
		INVALID_DATA: ResponseCode{
			code:       ovostatus.INVALID_DATA,
			message:    "Bank Account Not Found",
			httpStatus: http.StatusBadRequest,
		},
		THIRD_PARTY_ERROR: ResponseCode{
			code:       ovostatus.THIRD_PARTY_ERROR,
			message:    "Bank Account Not Found",
			httpStatus: http.StatusBadRequest,
		},
		INELIGIBLE: ResponseCode{
			code:       ovostatus.INELIGIBLE,
			message:    "Bank Account Not Found",
			httpStatus: http.StatusBadRequest,
		},
		PAYMENT_FAILED: ResponseCode{
			code:       ovostatus.PAYMENT_FAILED,
			message:    "Bank Account Not Found",
			httpStatus: http.StatusBadRequest,
		},
		PAYMENT_TIMED_OUT: ResponseCode{
			code:       ovostatus.PAYMENT_TIMED_OUT,
			message:    "Bank Account Not Found",
			httpStatus: http.StatusBadRequest,
		},
		INSUFFICIENT_FUND: ResponseCode{
			code:       ovostatus.INSUFFICIENT_FUND,
			message:    "Bank Account Not Found",
			httpStatus: http.StatusBadRequest,
		},
		INVALID_TRANSACTION: ResponseCode{
			code:       ovostatus.INVALID_TRANSACTION,
			message:    "Bank Account Not Found",
			httpStatus: http.StatusBadRequest,
		},
		ORDER_CREATION_FAILED: ResponseCode{
			code:       ovostatus.ORDER_CREATION_FAILED,
			message:    "Bank Account Not Found",
			httpStatus: http.StatusServiceUnavailable,
		},
		INCOMPLETE_ACCOUNT: ResponseCode{
			code:       ovostatus.INCOMPLETE_ACCOUNT,
			message:    "Account Is Not Complete",
			httpStatus: http.StatusUnprocessableEntity,
		},
		PROMO_CODE_ERROR: ResponseCode{
			code:       ovostatus.PROMO_CODE_ERROR,
			message:    "Promo Code Error",
			httpStatus: http.StatusBadRequest,
		},
		PROMO_BOOK_ERROR: ResponseCode{
			code:       ovostatus.PROMO_BOOK_ERROR,
			message:    "Promo Book Error",
			httpStatus: http.StatusBadRequest,
		},
		INSUFFICIENT_APP_VERSION: ResponseCode{
			code:       ovostatus.INSUFFICIENT_APP_VERSION,
			message:    "Insufficient App version",
			httpStatus: http.StatusBadRequest,
		},
		PENDING_TRANSACTION_ERROR: ResponseCode{
			code:       ovostatus.PENDING_TRANSACTION_ERROR,
			message:    "Failed",
			httpStatus: http.StatusBadRequest,
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
