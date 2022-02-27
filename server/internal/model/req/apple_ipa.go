package req

import "mime/multipart"

type IPAForm struct {
	IPA     *multipart.FileHeader `form:"ipa" binding:"required"`
	Summary string                `form:"summary"`
}

type IPAQuery struct {
	UUID string `form:"uuid" binding:"required"`
}
