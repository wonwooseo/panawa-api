package errors

import (
	"fmt"
	"net/http"
)

type RouterError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (re *RouterError) Error() string {
	return re.Message
}

func NewInternalServerError() (int, *RouterError) {
	return http.StatusInternalServerError, &RouterError{Code: "0000", Message: "internal server error"}
}

func NewInvalidQueryParamError(k, v string) (int, *RouterError) {
	return http.StatusBadRequest, &RouterError{Code: "1000", Message: fmt.Sprintf("invalid value %q for query param %q", v, k)}
}

func NewUnknownItemError(code string) (int, *RouterError) {
	return http.StatusBadRequest, &RouterError{Code: "1001", Message: fmt.Sprintf("unknown item code: %s", code)}
}

func NewUnknownRegionError(code string) (int, *RouterError) {
	return http.StatusBadRequest, &RouterError{Code: "1002", Message: fmt.Sprintf("unknown region code: %s", code)}
}

func NewNoPriceDataError() (int, *RouterError) {
	return http.StatusNotFound, &RouterError{Code: "2000", Message: "price data not found"}
}
