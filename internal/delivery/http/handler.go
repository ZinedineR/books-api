package http

import (
	"books-api/internal/delivery/http/response"
	"books-api/internal/model"
	"books-api/pkg/exception"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	invalidParameter = "invalid %s parameter"
)

const (
	filtersParam = "filter"
	orderParam   = "sort"
	pageParam    = "page"
	limitParam   = "pageSize"
)

var orderRegex = regexp.MustCompile("(\\w+):(\\w+)")

var OrderOperators = map[string]string{
	"desc": "desc",
	"asc":  "asc",
}

func GetOrderValue(value string) (string, error) {
	if op, ok := OrderOperators[value]; ok {
		return op, nil
	}
	return "", fmt.Errorf(invalidParameter, value)
}

var filterRegex = regexp.MustCompile(`(\w+):([^|]+):(\w+)`)

var FilterOperator = map[string]string{
	"eq":   "=",
	"lt":   "<",
	"gt":   ">",
	"lte":  "<=",
	"gte":  ">=",
	"in":   "in",
	"like": "like",
	"is":   "is",
	"not":  "not in",
}

func GetFilterOperator(operator string) (string, error) {
	if op, ok := FilterOperator[operator]; ok {
		return op, nil
	}
	return "", fmt.Errorf(invalidParameter, operator)
}

type Handler struct {
}

func (h *Handler) JSON(e *gin.Context, r response.IResponse) {
	e.JSON(r.GetStatusCode(), r)
}

func (h *Handler) AbortJSON(e *gin.Context, r response.IResponse) {
	e.AbortWithStatusJSON(r.GetStatusCode(), r)
}
func (h *Handler) InternalErrorJSON(e *gin.Context, msg any, err ...any) {
	h.ErrorJSON(e, 500, msg, err)
}

func (h *Handler) PaginationJSON(c *gin.Context, pagination any, data any) {
	h.JSON(c, &response.PaginationResponse{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success",
		Pagination:      pagination,
		Data:            data,
	})
}

func (h *Handler) SuccessJSON(c *gin.Context) {
	h.JSON(c, &response.SuccessResponse{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success",
	})
}

func (h *Handler) SuccessMessageJSON(c *gin.Context, message string) {
	h.JSON(c, &response.SuccessResponse{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success",
	})
}

func (h *Handler) DataJSON(c *gin.Context, data any) {
	h.JSON(c, &response.DataResponse{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success",
		Data:            data,
	})
}
func (h *Handler) ExceptionJSON(e *gin.Context, exc *exception.Exception) {
	h.AbortJSON(e, &response.ErrorResponse{
		ResponseCode:    exc.GetHttpCode(),
		ResponseMessage: exc.Message,
		Error:           exc.GetError(),
	})
}

func (h *Handler) ErrorJSON(e *gin.Context, status int, msg any, err ...any) {
	var er any
	if len(err) > 0 {
		er = err[0]
	}
	h.AbortJSON(e, &response.ErrorResponse{
		ResponseCode:    status,
		ResponseMessage: msg,
		Error:           er,
	})
}

func (h *Handler) BadRequestJSON(e *gin.Context, msg any, err ...any) {
	h.ErrorJSON(e, 400, msg, err)
}

func (h *Handler) UnauthorizedJSON(e *gin.Context, msg any, err ...any) {
	h.ErrorJSON(e, 401, msg, err)
}

func (h *Handler) SignatureJSON(c *gin.Context, signature string) {
	h.JSON(c, &response.DataResponse{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success",
		Data:            &model.Signature{Signature: signature},
	})
}

func (h *Handler) ParamInt(e *gin.Context, key string) (int, error) {
	return strconv.Atoi(e.Param(key))
}

func (h *Handler) ParamInt64(e *gin.Context, key string) (int64, error) {
	return strconv.ParseInt(e.Param(key), 10, 64)
}

func (h *Handler) ParseHTTPMethod(c *gin.Context) (string, string, error) {
	httpMethod := c.GetHeader("httpMethod")

	// Validate HTTP method
	if httpMethod != http.MethodPost && httpMethod != http.MethodGet &&
		httpMethod != http.MethodPut && httpMethod != http.MethodDelete {
		return "", "", errors.New("http method invalid")
	}

	// Read the request body
	bodyJson, errRead := io.ReadAll(c.Request.Body)
	if errRead != nil {
		return "", "", errRead
	}

	// Check if the body is non-empty and valid JSON
	if len(bodyJson) != 0 {
		var x interface{}
		if err := json.Unmarshal(bodyJson, &x); err != nil {
			return "", "", errors.New("body is not valid JSON")
		}
	}

	// Return the HTTP method and body (as []byte)
	return httpMethod, string(bodyJson), nil
}

//	func (h *Handler) ParseSignatureHTTPMethod(c *gin.Context) (string, string, error) {
//		httpMethod := c.Request.Method
//
//		// Read the request body
//		//bodyJson, errRead := io.ReadAll(c.Request.Body)
//		//if errRead != nil {
//		//	return "", "", errRead
//		//}
//		//
//		//// Check if the body is non-empty and valid JSON
//		//if len(bodyJson) != 0 {
//		//	var x interface{}
//		//	if err := json.Unmarshal(bodyJson, &x); err != nil {
//		//		return "", "", errors.New("body is not valid JSON")
//		//	}
//		//}
//		var bodyBytes []byte
//		if c.Request.Body != nil {
//			bodyBytes, _ = io.ReadAll(c.Request.Body)
//		}
//		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
//
//		if (c.Request.Method == "POST" || c.Request.Method == "PUT") && len(bodyBytes) == 0 {
//			return "", "", errors.New("body is required")
//		}
//		// Return the HTTP method and body (as []byte)
//		return httpMethod, string(bodyBytes), nil
//	}
func (h *Handler) ParseSignatureHTTPMethod(c *gin.Context) (string, string, error) {
	httpMethod := c.Request.Method

	// Read the request body
	bodyBytes, errRead := io.ReadAll(c.Request.Body)
	if errRead != nil {
		return "", "", errRead
	}

	// Restore the body so it can be read again later
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if (c.Request.Method == "POST" || c.Request.Method == "PUT") && len(bodyBytes) == 0 {
		return "", "", errors.New("body is required")
	}

	// Return the HTTP method and body (as string)
	return httpMethod, string(bodyBytes), nil
}
func (h *Handler) GetToken(c *gin.Context) string {
	return c.GetString("access_token")
}

func (h *Handler) ParseNameParam(c *gin.Context) (string, string) {
	nameQuery := c.Query("name")
	if nameQuery == "" {
		return "", ""
	}
	filterQuery := strings.Split(nameQuery, ":")
	return filterQuery[0], filterQuery[1]
}

func (h *Handler) ParseDateParam(c *gin.Context) (time.Time, time.Time, error) {
	from := c.Query("from")
	to := c.Query("to")
	toDate := time.Now().AddDate(0, 0, 1)
	var fromDate time.Time
	var err error
	if from != "" {
		fromDate, err = time.Parse("2006-01-02", from)
		if err != nil {
			return fromDate, toDate, err
		}
	}
	if to != "" {
		toDate, err = time.Parse("2006-01-02", to)
		if err != nil {
			return fromDate, toDate, err
		}
		toDate = toDate.AddDate(0, 0, 1)
	}
	return fromDate, toDate, nil
}
func (h *Handler) ParsePageParam(c *gin.Context) (int64, int64, error) {
	var err error
	limit, err := strconv.Atoi(c.DefaultQuery(limitParam, "0"))
	page, err := strconv.Atoi(c.DefaultQuery(pageParam, "0"))
	if err != nil {
		return 0, 0, err
	}
	return int64(limit), int64(page), nil
}

func (h *Handler) ParsePageLimitParam(c *gin.Context) (model.PaginationParam, error) {
	var p model.PaginationParam
	var err error
	p.Page, err = strconv.Atoi(c.DefaultQuery(pageParam, "1"))
	p.PageSize, err = strconv.Atoi(c.DefaultQuery(limitParam, "10"))
	if err != nil {
		return model.PaginationParam{}, err
	}
	return p, nil
}

func (h *Handler) ParseOrderParam(c *gin.Context) (model.OrderParam, error) {
	var p model.OrderParam
	order := c.Query(orderParam)
	if order != "" {
		listOrder := strings.Split(order, ",")
		for _, o := range listOrder {
			if !orderRegex.MatchString(o) {
				continue
			}
			condition := strings.Split(o, ":")
			if len(condition) != 2 {
				return model.OrderParam{}, fmt.Errorf(invalidParameter, orderParam)
			}
			value, err := GetOrderValue(condition[1])
			if err != nil {
				return model.OrderParam{}, err
			}
			p.OrderBy = condition[0]
			p.Order = value
		}
	}
	return p, nil
}

func (h *Handler) ParseFilterParams(c *gin.Context) (model.FilterParams, error) {
	var p model.FilterParams
	f := c.Query(filtersParam)

	if f != "" {
		listFilter := strings.Split(f, "|")
		for _, v := range listFilter {
			if !filterRegex.MatchString(v) {
				continue
			}
			filter := strings.Split(v, ":")
			if len(filter) != 3 {
				return model.FilterParams{}, fmt.Errorf(invalidParameter, filter)
			}
			operator, err := GetFilterOperator(filter[2])
			if err != nil {
				return model.FilterParams{}, err
			}
			p = append(p, &model.FilterParam{
				Field:    filter[0],
				Value:    filter[1],
				Operator: operator,
			})
		}
	}

	return p, nil
}

func (h *Handler) ParsePaginationParams(c *gin.Context) (
	model.PaginationParam, model.OrderParam, model.FilterParams, error,
) {
	page, err := h.ParsePageLimitParam(c)
	if err != nil {
		return model.PaginationParam{}, model.OrderParam{}, model.FilterParams{}, err
	}
	order, err := h.ParseOrderParam(c)
	if err != nil {
		return model.PaginationParam{}, model.OrderParam{}, model.FilterParams{}, err
	}
	filters, err := h.ParseFilterParams(c)
	if err != nil {
		return model.PaginationParam{}, model.OrderParam{}, model.FilterParams{}, err
	}
	return page, order, filters, nil
}
