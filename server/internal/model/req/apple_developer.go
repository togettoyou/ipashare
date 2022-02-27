package req

import "mime/multipart"

type AppleDeveloperForm struct {
	P8  *multipart.FileHeader `form:"p8" binding:"required"`
	Iss string                `form:"iss" binding:"required"`
	Kid string                `form:"kid" binding:"required"`
}

type AppleDeveloperQuery struct {
	Iss string `form:"iss" binding:"required"`
}
