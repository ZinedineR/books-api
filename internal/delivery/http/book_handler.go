package http

import (
	_ "books-api/internal/delivery/http/response"
	"books-api/internal/entity"
	"books-api/internal/model"
	service "books-api/internal/services"
	"github.com/gin-gonic/gin"
)

type BookHTTPHandler struct {
	Handler
	BookService service.BookService
}

func NewBookHTTPHandler(book service.BookService) *BookHTTPHandler {
	return &BookHTTPHandler{
		BookService: book,
	}
}

// Create godoc
// @Summary Create a new book
// @Description Creates a new book with the provided details
// @Tags Books
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param X-Req-Signature header string true "generated from signature"
// @Param notification-list body entity.UpsertBook true "Book Request"
// @Success 200 {object} response.DataResponse{data=entity.UpsertBook} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /books [post]
func (h BookHTTPHandler) Create(ctx *gin.Context) {
	request := entity.UpsertBook{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	if errException := h.BookService.Create(ctx, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

// List godoc
// @Summary List books
// @Description Retrieves a paginated list of books with optional ordering and filtering
// @Tags Books
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param X-Req-Signature header string true "generated from signature"
// @Param pageSize query string false "Number of items per page" default(0)
// @Param page query string false "Page number" default(0)
// @Param filter query string false "Filter rules<br><br>### Rules Filter<br>rule:<br>  * {Name of Field}:{value}:{Symbol}<br><br>Symbols:<br>  * eq (=)<br>  * lt (<)<br>  * gt (>)<br>  * lte (<=)<br>  * gte (>=)<br>  * in (in)<br>  * like (like)<br><br>Field list:<br>  * id<br>  * title<br>  * isbn<br>  * author_id" default(id:1:eq)
// @Param sort query string false "Sort rules:<br><br>### Rules Sort<br>rule:<br>  * {Name of Field}:{Symbol}<br><br>Symbols:<br>  * asc<br>  * desc<br><br>Field list:<br>  * id<br>  * title<br>  * isbn<br>  * author_id" default(id:desc)
// @Success 200 {object} response.PaginationResponse{data=[]entity.Book,pagination=model.Pagination} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /books [get]
func (h BookHTTPHandler) List(ctx *gin.Context) {
	var req model.ListReq
	var err error
	req.Page, req.Order, req.Filter, err = h.ParsePaginationParams(ctx)
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	result, errException := h.BookService.List(ctx, req)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

// FindOne godoc
// @Summary Get details of a book
// @Description Retrieves the details of a specific book by ID
// @Tags Books
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param X-Req-Signature header string true "generated from signature"
// @Param id path string true "Book ID (UUID format)"
// @Success 200 {object} response.DataResponse{data=entity.Book} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /books/{id} [get]
func (h BookHTTPHandler) FindOne(ctx *gin.Context) {
	idParam := ctx.Param("id")
	result, errException := h.BookService.FindOne(ctx, idParam)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}

// Update godoc
// @Summary Update an existing book
// @Description Updates an existing book with the provided details
// @Tags Books
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param X-Req-Signature header string true "generated from signature"
// @Param id path string true "Book ID (UUID format)"
// @Param book body entity.UpsertBook true "Updated Book details"
// @Success 200 {object} response.DataResponse{data=entity.UpsertBook} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /books/{id} [put]
func (h BookHTTPHandler) Update(ctx *gin.Context) {
	// Get Info
	idParam := ctx.Param("id")
	request := entity.UpsertBook{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	if errException := h.BookService.Update(ctx, idParam, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

// Delete godoc
// @Summary Delete an existing book
// @Description Deletes an existing book by ID
// @Tags Books
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param X-Req-Signature header string true "generated from signature"
// @Param id path string true "Book ID (UUID format)"
// @Success 200 {object} response.SuccessResponse "success"
// @Failure 400 {object} response.SuccessResponse "error"
// @Router /books/{id} [delete]
func (h BookHTTPHandler) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if errException := h.BookService.Delete(ctx, idParam); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.SuccessMessageJSON(ctx, idParam+" has been deleted")
}
