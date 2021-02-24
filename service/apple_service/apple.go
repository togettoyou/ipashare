package apple_service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"super-signature/models"
	"super-signature/pkg/apple"
	"super-signature/pkg/setting"
	"super-signature/pkg/util"
)

func AddAppleAccount(iss string, kid string, p8 string, csr string) (int, error) {
	// 判断数据库中是否已经存在
	appleAccount, err := models.GetAppleAccountByIss(iss)
	if err != nil {
		return 0, err
	}
	if appleAccount != nil {
		return 0, errors.New("开发者账号已存在,不能重复添加")
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
		return 0, errors.New("可用设备已不足")
	}
	// 删除账号下所有的Certificates和BundleIds
	//err = authorize.DeleteAllBundleIds()
	//if err != nil {
	//	return 0, err
	//}
	//err = authorize.DeleteAllCertificates()
	//if err != nil {
	//	return 0, err
	//}
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
	var cerPath = setting.PathSetting.AppleAccountPath + iss + "/cer.cer"
	err = util.Base64ToFile(certificateResponse.Data.Attributes.CertificateContent, cerPath)
	if err != nil {
		return 0, err
	}
	// 根据cerFile调用openssl生成pem
	var pemPath = fmt.Sprintf("%s%s/pem.pem", setting.PathSetting.AppleAccountPath, iss)
	err = util.RunCmd(
		fmt.Sprintf(
			"openssl x509 -in %s -inform DER -outform PEM -out %s",
			cerPath,
			pemPath),
	)
	if err != nil {
		log.Printf("%s", err.Error())
		return 0, err
	}
	// 将账户信息插入到数据库
	account := models.AppleAccount{
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
		appleDevice := models.AppleDevice{
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

//删除指定开发者账号
func DeleteAppleAccountByIss(iss string) error {
	appleAccount, err := models.GetAppleAccountByIss(iss)
	if err != nil {
		return err
	}
	if appleAccount == nil {
		return errors.New("开发者账号不存在")
	}
	log.Println("------开始删除开发者账号------")
	err = models.DeleteAppleAccountByIss(iss)
	if err != nil {
		return err
	}
	//清除对应设备
	log.Println("清除对应设备...")
	err = models.DeleteAppleDeviceByAccountIss(iss)
	if err != nil {
		return err
	}
	//删除cer
	log.Println("删除cer...")
	err = os.Remove(appleAccount.CerPath)
	if err != nil {
		return err
	}
	//删除pem
	log.Println("删除pem...")
	err = os.Remove(appleAccount.PemPath)
	if err != nil {
		return err
	}
	//清除开发者账户中心cer证书
	log.Println("清除开发者账户中心cer证书...")
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
