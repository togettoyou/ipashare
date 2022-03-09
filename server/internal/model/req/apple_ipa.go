package req

import "mime/multipart"

type IPAForm struct {
	IPA     *multipart.FileHeader `form:"ipa" binding:"required"`
	Summary string                `form:"summary" binding:"required,min=2,max=100"`
}

type IPAQuery struct {
	UUID string `form:"uuid" binding:"required"`
}

type IPABody struct {
	UUID    string `json:"uuid" binding:"required"`
	Summary string `json:"summary" binding:"required,min=2,max=100"`
}
