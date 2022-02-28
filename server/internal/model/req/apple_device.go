package req

type AppleDeviceUri struct {
	UUID string `uri:"uuid" binding:"required"`
}
