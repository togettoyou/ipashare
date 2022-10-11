package req

type AppleDeviceList struct {
	Iss string `form:"iss" binding:"required"`
}

type AppleDeviceUri struct {
	UUID string `uri:"uuid" binding:"required"`
}
