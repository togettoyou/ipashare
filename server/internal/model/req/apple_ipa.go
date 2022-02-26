package req

import "mime/multipart"

type IPAForm struct {
	IPA     *multipart.FileHeader `form:"ipa" binding:"required"`
	Summary string                `form:"summary"`
}
