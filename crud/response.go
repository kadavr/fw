package crud

import (
	"fmt"

	"gorm.io/gorm"
)

const (
	ErrorUnknown             = 0
	ErrorNotFound            = 1
	ErrorInvalidRequest      = 2
	ErrorAuthNoSessionCookie = 3
	ErrorAuthServiceError    = 4
	ErrorAuthNoPhone         = 5
	ErrorAuthWrongPhone      = 6
	ErrorTelegramSendError   = 7
	ErrorTelegramNoBotToken  = 8
	ErrorTelegramNoChatID    = 9
	ErrorTooManyResult       = 10
	ErrorShortRequest        = 11
	ErrorEmployeeIsNull      = 12
	ErrorTimeout             = 14
)

// Responser ...
type Responser[T Model] interface {
	ResponseOne(code int64, body T) response[T]
	ResponseMany(code int64, rangeInfo rangeConf, total int64, body []T) listResponse[T]
}

// ErrorResponser ...
type ErrorResponser interface {
	Match(error) errorResponse
}

type errorResponse struct {
	Status  int64  `json:"status"`
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// ListRange описывает пейджер в ответе
type listRange struct {
	Count  int64 `json:"count"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

// ListResponse ответ со списком элементов
type listResponse[T any] struct {
	Status int64       `json:"status"`
	Data   listData[T] `json:"data"`
}

// ListData список элементов ответа
type listData[T any] struct {
	Items []T       `json:"items"`
	Range listRange `json:"range"`
}

// Response ответ с одним элементом
type response[T Model] struct {
	Status int64 `json:"status"`
	Data   T     `json:"data"`
}

// NewResponse ...
func NewResponse[T Model](status int64, body T) response[T] {
	return response[T]{
		Status: status,
		Data:   body,
	}
}

type rangeConf struct {
	Limit  int64
	Offset int64
}

// NewListResponse ...
func NewListResponse[T any](status int64, rangeInfo rangeConf, total int64, body []T) listResponse[T] {
	return listResponse[T]{
		Status: status,
		Data: listData[T]{
			Items: body,
			Range: listRange{
				Count:  total,
				Limit:  rangeInfo.Limit,
				Offset: rangeInfo.Offset,
			},
		},
	}
}

type respOk[T Model] struct {
	// resps map[OK]
}

// NewResponseRegistry ...
func NewResponseRegistry[T Model]() respOk[T] {
	return respOk[T]{}
}

func (r respOk[T]) ResponseOne(code int64, body T) response[T] {
	return NewResponse(code, body)
}
func (r respOk[T]) ResponseMany(code int64, rangeInfo rangeConf, total int64, body []T) listResponse[T] {
	return NewListResponse(code, rangeInfo, total, body)
}

type errRepo struct {
	errors map[error]errorResponse
}

// NewErrResponseRegistry ...
func NewErrResponseRegistry(err ...error) errRepo {
	return errRepo{
		errors: map[error]errorResponse{
			gorm.ErrRecordNotFound: {
				Status:  400,
				Code:    ErrorNotFound,
				Message: gorm.ErrRecordNotFound.Error(),
			},
		},
	}
}

// Match ...
func (r errRepo) Match(err error) errorResponse {
	fmt.Println(err)
	return r.errors[err]
}
