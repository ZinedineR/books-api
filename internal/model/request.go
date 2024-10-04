package model

type ListReq struct {
	Page   PaginationParam
	Order  OrderParam
	Filter FilterParams
}

type UpdateApproval struct {
	Id      []string `json:"id"`
	Remarks string   `json:"remarks"`
}
