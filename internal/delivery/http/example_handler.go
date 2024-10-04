package http

import (
	_ "boiler-plate-clean/internal/delivery/http/response"
	"boiler-plate-clean/internal/entity"
	service "boiler-plate-clean/internal/services"
	"github.com/gin-gonic/gin"
)

type ExampleHTTPHandler struct {
	Handler
	ExampleService service.ExampleService
}

func NewExampleHTTPHandler(example service.ExampleService) *ExampleHTTPHandler {
	return &ExampleHTTPHandler{
		ExampleService: example,
	}
}

func (h ExampleHTTPHandler) Create(ctx *gin.Context) {
	request := entity.Example{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	if errException := h.ExampleService.CreateExample(ctx, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

func (h ExampleHTTPHandler) Find(ctx *gin.Context) {
	//var req model.ListReq
	//var err error
	//req.Page, req.Order, req.Filter, err = h.ParsePaginationParams(ctx)
	//if err != nil {
	//	h.BadRequestJSON(ctx, err.Error())
	//	return
	//}
	//result, errException := h.ExampleService.Find(ctx, req)
	//if errException != nil {
	//	h.ExceptionJSON(ctx, errException)
	//	return
	//}

	//h.DataJSON(ctx, result)
}

func (h ExampleHTTPHandler) FindOne(ctx *gin.Context) {
	//idParam := ctx.Param("id")
	//result, errException := h.ExampleService.FindOne(ctx, idParam)
	//if errException != nil {
	//	h.ExceptionJSON(ctx, errException)
	//	return
	//}

	//h.DataJSON(ctx, result)
}

func (h ExampleHTTPHandler) Update(ctx *gin.Context) {
	//// Get Info
	////idParam := ctx.Param("id")
	//request := entity.UpsertExample{}
	//if err := ctx.ShouldBindJSON(&request); err != nil {
	//	h.BadRequestJSON(ctx, err.Error())
	//	return
	//}
	//
	////if errException := h.ExampleService.Update(ctx, idParam, &request, requestHeader); errException != nil {
	////	h.ExceptionJSON(ctx, errException)
	////	return
	////}
	//
	//h.DataJSON(ctx, request)
}

func (h ExampleHTTPHandler) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	//if errException := h.ExampleService.Delete(ctx, idParam, requestHeader); errException != nil {
	//	h.ExceptionJSON(ctx, errException)
	//	return
	//}

	h.SuccessMessageJSON(ctx, idParam+" has been deleted")
}
