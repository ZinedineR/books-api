package model

type FilterParam struct {
	Field    string
	Value    string
	Operator string
}

type FilterParams []*FilterParam
