package http

import (
	_ "books-api/internal/delivery/http/response"
	_ "books-api/internal/model"
	"books-api/pkg/signature"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type SignatureHTTPHandler struct {
	Handler
	Signaturer signature.Signaturer
}

func NewSignatureHTTPHandler(signaturer signature.Signaturer) *SignatureHTTPHandler {
	return &SignatureHTTPHandler{
		Signaturer: signaturer,
	}
}

// Signature godoc
// @Summary Signature API Request
// @Description Generate Signature to authenticate api request
// @Tags Signatures
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param httpMethod header string true "HTTP Method"
// @Param body body object false "Optional JSON body"
// @Success 200 {object} response.DataResponse{data=model.Signature} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /auth/signature [post]
func (h SignatureHTTPHandler) Signature(ctx *gin.Context) {
	var err error
	httpMethod, bodyJson, err := h.ParseHTTPMethod(ctx)
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	token := h.GetToken(ctx)
	signatureResult, err := h.Signaturer.SignHMAC512(httpMethod, bodyJson, token)
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	slog.Info("signatureResult", slog.String("signatureResult", signatureResult))
	h.SignatureJSON(ctx, signatureResult)
}
