package apple_service

import (
	"fmt"
	"os"
	"super-signature/model"
	"super-signature/util/apple"
	"super-signature/util/conf"
	"super-signature/util/errno"
	"super-signature/util/tools"
)

func AddAppleAccount(iss string, kid string, p8 string, csr string) (int, error) {
	appleAccount, err := model.GetAppleAccountByIss(iss)
	if err != nil {
		return 0, err
	}
	if appleAccount != nil {
		return 0, errno.ErrHaveAppleAccount
	}
	var authorize = apple.Authorize{
		P8:  p8,
		Iss: iss,
		Kid: kid,
	}
	// 验证账户合法性
	devicesResponseList, err := authorize.GetAvailableDevices()
	if err != nil {
		return 0, err
	}
	// 判断账户可用设备是否充足
	if devicesResponseList.Meta.Paging.Total >= 100 {
		return 0, errno.ErrDeviceInsufficient
	}
	if conf.Config.AppleConf.DeleteAllCertificates {
		// 删除账号下所有的Certificates
		err = authorize.DeleteAllCertificates()
		if err != nil {
			return 0, err
		}
	}
	// 判断账号是否存在bundleIds为*的数据
	bundleIds, err := authorize.GetBundleIdsByIdentifier("*")
	if err != nil {
		return 0, err
	}
	if bundleIds == "" {
		// 创建新的bundleIds
		bundleIdResponse, err := authorize.CreateBundleIds("*")
		if err != nil {
			return 0, err
		}
		bundleIds = bundleIdResponse.Data.ID
	}
	// 根据csr创建新的CerId和cer文件证书
	certificateResponse, err := authorize.CreateCertificates(csr)
	if err != nil {
		return 0, err
	}
	var cerId = certificateResponse.Data.ID
	var cerPath = conf.Config.ApplePath.AppleAccountPath + iss + "/cer.cer"
	err = tools.Base64ToFile(certificateResponse.Data.Attributes.CertificateContent, cerPath)
	if err != nil {
		return 0, err
	}
	// 根据cerFile调用openssl生成pem
	var pemPath = fmt.Sprintf("%s%s/pem.pem", conf.Config.ApplePath.AppleAccountPath, iss)
	err = tools.RunCmd(
		fmt.Sprintf(
			"openssl x509 -in %s -inform DER -outform PEM -out %s",
			cerPath,
			pemPath),
	)
	if err != nil {
		return 0, err
	}
	// 将账户信息插入到数据库
	account := model.AppleAccount{
		Iss:       iss,
		Kid:       kid,
		CerId:     cerId,
		P8:        p8,
		CerPath:   cerPath,
		PemPath:   pemPath,
		BundleIds: bundleIds,
		Count:     devicesResponseList.Meta.Paging.Total,
	}
	err = account.InsertAppleAccount()
	if err != nil {
		return 0, err
	}
	// 将当前账户存在的udid添加到数据库
	for _, v := range devicesResponseList.Data {
		appleDevice := model.AppleDevice{
			AccountIss: iss,
			Udid:       v.Attributes.Udid,
			DeviceId:   v.ID,
		}
		if err = appleDevice.InsertAppleDevice(); err != nil {
			return 0, err
		}
	}
	return 100 - devicesResponseList.Meta.Paging.Total, nil
}

// DeleteAppleAccountByIss 删除指定开发者账号
func DeleteAppleAccountByIss(iss string) error {
	appleAccount, err := model.GetAppleAccountByIss(iss)
	if err != nil {
		return err
	}
	if appleAccount == nil {
		return errno.ErrNotAppleAccount
	}
	err = model.DeleteAppleAccountByIss(iss)
	if err != nil {
		return err
	}
	//清除对应设备
	err = model.DeleteAppleDeviceByAccountIss(iss)
	if err != nil {
		return err
	}
	//删除cer
	err = os.Remove(appleAccount.CerPath)
	if err != nil {
		return err
	}
	//删除pem
	err = os.Remove(appleAccount.PemPath)
	if err != nil {
		return err
	}
	//清除开发者账户中心cer证书
	err = apple.Authorize{
		P8:  appleAccount.P8,
		Iss: appleAccount.Iss,
		Kid: appleAccount.Kid,
	}.DeleteCertificatesByCerId(appleAccount.CerId)
	if err != nil {
		return err
	}
	return nil
}
