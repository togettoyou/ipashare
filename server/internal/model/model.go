package model

// Store 实体管理，所有DB操作
type Store struct {
	AppleDeveloper AppleDeveloperStore
	AppleDevice    AppleDeviceStore
	AppleIPA       AppleIPAStore
}
