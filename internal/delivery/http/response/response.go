package response

type IResponse interface {
	GetStatusCode() int
}
type Status struct {
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}
type ErrorResponse struct {
	ResponseCode    int `json:"responseCode"`
	ResponseMessage any `json:"responseMessage"`
	Error           any `json:"error"`
}

func (r *ErrorResponse) GetStatusCode() int {
	return r.ResponseCode
}

type SuccessResponse struct {
	ResponseCode    int `json:"responseCode"`
	ResponseMessage any `json:"responseMessage"`
}

func (r *SuccessResponse) GetStatusCode() int {
	return r.ResponseCode
}

type DataResponse struct {
	ResponseCode    int `json:"responseCode"`
	ResponseMessage any `json:"responseMessage"`
	Data            any `json:"data"`
}

func (r *DataResponse) GetStatusCode() int {
	return r.ResponseCode
}

type PaginationResponse struct {
	ResponseCode    int `json:"responseCode"`
	ResponseMessage any `json:"responseMessage"`
	Pagination      any `json:"pagination"`
	Data            any `json:"data"`
}

func (r *PaginationResponse) GetStatusCode() int {
	return r.ResponseCode
}
