package http

import (
	_ "books-api/internal/delivery/http/response"
	"books-api/internal/entity"
	"books-api/internal/model"
	service "books-api/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthorHTTPHandler struct {
	Handler
	AuthorService service.AuthorService
}

func NewAuthorHTTPHandler(author service.AuthorService) *AuthorHTTPHandler {
	return &AuthorHTTPHandler{
		AuthorService: author,
	}
}

// Create godoc
// @Summary Create a new author
// @Description Creates a new author with the provided details
// @Tags Authors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param X-Req-Signature header string true "generated from signature"
// @Param author body entity.UpsertAuthor true "Author Request"
// @Success 200 {object} response.DataResponse{data=entity.UpsertAuthor} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /authors [post]
func (h AuthorHTTPHandler) Create(ctx *gin.Context) {
	request := entity.UpsertAuthor{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	if errException := h.AuthorService.Create(ctx, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

// List godoc
// @Summary List authors
// @Description Retrieves a paginated list of authors with optional ordering and filtering
// @Tags Authors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param X-Req-Signature header string true "generated from signature"
// @Param pageSize query string false "Number of items per page" default(0)
// @Param page query string false "Page number" default(0)
// @Param filter query string false "Filter rules<br><br>### Rules Filter<br>rule:<br>  * {Name of Field}:{value}:{Symbol}<br><br>Symbols:<br>  * eq (=)<br>  * lt (<)<br>  * gt (>)<br>  * lte (<=)<br>  * gte (>=)<br>  * in (in)<br>  * like (like)<br><br>Field list:<br>  * id<br>  * name<br>  * birthdate" default(id:1:eq)
// @Param sort query string false "Sort rules:<br><br>### Rules Sort<br>rule:<br>  * {Name of Field}:{Symbol}<br><br>Symbols:<br>  * asc<br>  * desc<br><br>Field list:<br>  * id<br>  * name<br>  * birthdate" default(id:desc)
// @Success 200 {object} response.PaginationResponse{data=[]entity.Author,pagination=model.Pagination} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /authors [get]
func (h AuthorHTTPHandler) List(ctx *gin.Context) {
	var req model.ListReq
	var err error
	req.Page, req.Order, req.Filter, err = h.ParsePaginationParams(ctx)
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	result, errException := h.AuthorService.List(ctx, req)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

// FindOne godoc
// @Summary Get details of an author
// @Description Retrieves the details of a specific author by ID
// @Tags Authors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param X-Req-Signature header string true "generated from signature"
// @Param id path string true "Author ID (UUID format)"
// @Success 200 {object} response.DataResponse{data=entity.Author} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /authors/{id} [get]
func (h AuthorHTTPHandler) FindOne(ctx *gin.Context) {
	idParam := ctx.Param("id")
	result, errException := h.AuthorService.FindOne(ctx, idParam)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

// Update godoc
// @Summary Update an existing author
// @Description Updates an existing author with the provided details
// @Tags Authors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param X-Req-Signature header string true "generated from signature"
// @Param id path string true "Author ID (UUID format)"
// @Param author body entity.UpsertAuthor true "Updated Author details"
// @Success 200 {object} response.DataResponse{data=entity.UpsertAuthor} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /authors/{id} [put]
func (h AuthorHTTPHandler) Update(ctx *gin.Context) {
	// Get Info
	idParam := ctx.Param("id")
	request := entity.UpsertAuthor{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	if errException := h.AuthorService.Update(ctx, idParam, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

// Delete godoc
// @Summary Delete an existing author
// @Description Deletes an existing author by ID
// @Tags Authors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param X-Req-Signature header string true "generated from signature"
// @Param id path string true "Author ID (UUID format)"
// @Success 200 {object} response.SuccessResponse "success"
// @Failure 400 {object} response.SuccessResponse "error"
// @Router /authors/{id} [delete]
func (h AuthorHTTPHandler) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if errException := h.AuthorService.Delete(ctx, idParam); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.SuccessMessageJSON(ctx, idParam+" has been deleted")
}
