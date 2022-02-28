package svc

import (
	"io/ioutil"
	"os"
	"path"
	"supersign/internal/model"
	"supersign/pkg/appstore"
	"supersign/pkg/conf"
	"supersign/pkg/e"
	"supersign/pkg/openssl"
	"supersign/pkg/tools"

	"gorm.io/gorm"
)

type AppleDeveloper struct {
	Service
}

func (a *AppleDeveloper) Add(iss, kid, p8 string) (num int, err error) {
	defer func() {
		if err != nil {
			os.RemoveAll(path.Join(conf.Apple.AppleDeveloperPath, iss))
		}
	}()
	appleDeveloper, err := a.store.AppleDeveloper.Query(iss)
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, e.NewWithStack(e.DBError, err)
	}
	if appleDeveloper != nil {
		return 0, e.ErrIssExist
	}
	authorize := appstore.Authorize{
		P8:  p8,
		Iss: iss,
		Kid: kid,
	}
	// 验证账号合法性
	devices, err := authorize.GetAvailableDevices()
	if err != nil {
		return 0, e.NewWithStack(e.ErrAppstoreAPI, err)
	}
	// 判断账户可用设备是否充足
	if devices.Meta.Paging.Total >= 100 {
		return 0, e.ErrDeviceInsufficient
	}
	// 生成CSR和KEY证书
	csrPath := path.Join(conf.Apple.AppleDeveloperPath, iss, "csr.csr")
	keyPath := path.Join(conf.Apple.AppleDeveloperPath, iss, "key.key")
	err = openssl.GenKeyAndReqCSR(keyPath, csrPath)
	if err != nil {
		return 0, e.NewWithStack(e.ErrIssAdd, err)
	}
	// 判断账号是否存在bundleIds为*的记录
	bundleIds, err := authorize.GetBundleIdsByIdentifier("*")
	if err != nil {
		return 0, e.NewWithStack(e.ErrAppstoreAPI, err)
	}
	if bundleIds == "" {
		// 创建新的bundleIds
		bundleIdResponse, err := authorize.CreateBundleIds("*")
		if err != nil {
			return 0, e.NewWithStack(e.ErrAppstoreAPI, err)
		}
		bundleIds = bundleIdResponse.Data.ID
	}
	// 根据CSR创建新的CerId和cer文件证书
	csrFile, err := os.Open(csrPath)
	defer csrFile.Close()
	if err != nil {
		return 0, e.NewWithStack(e.ErrIssAdd, err)
	}
	csrBytes, err := ioutil.ReadAll(csrFile)
	if err != nil {
		return 0, e.NewWithStack(e.ErrIssAdd, err)
	}
	certificateResponse, err := authorize.CreateCertificates(string(csrBytes))
	if err != nil {
		return 0, e.NewWithStack(e.ErrAppstoreAPI, err)
	}
	cerPath := path.Join(conf.Apple.AppleDeveloperPath, iss, "cer.cer")
	err = tools.Base64ToFile(certificateResponse.Data.Attributes.CertificateContent, cerPath)
	if err != nil {
		return 0, e.NewWithStack(e.ErrIssAdd, err)
	}
	// 根据cer生成pem
	pemPath := path.Join(conf.Apple.AppleDeveloperPath, iss, "pem.pem")
	err = openssl.GenPEM(cerPath, pemPath)
	if err != nil {
		return 0, e.NewWithStack(e.ErrIssAdd, err)
	}
	// 插入数据库
	appleDevices := make([]model.AppleDevice, 0)
	for _, datum := range devices.Data {
		appleDevices = append(appleDevices, model.AppleDevice{
			UDID:     datum.Attributes.Udid,
			Iss:      iss,
			DeviceID: datum.ID,
		})
	}
	err = a.store.AppleDeveloper.Create(&model.AppleDeveloper{
		Iss:       iss,
		Kid:       kid,
		P8:        p8,
		BundleIds: bundleIds,
		CerID:     certificateResponse.Data.ID,
		Count:     devices.Meta.Paging.Total,
		Limit:     100,
		Enable:    true,
	}, appleDevices)
	if err != nil {
		return 0, e.NewWithStack(e.DBError, err)
	}
	return 100 - devices.Meta.Paging.Total, nil
}

func (a *AppleDeveloper) Del(iss string) error {
	appleDeveloper, err := a.store.AppleDeveloper.Query(iss)
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	err = a.store.AppleDeveloper.Del(iss)
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	os.RemoveAll(path.Join(conf.Apple.AppleDeveloperPath, iss))
	authorize := appstore.Authorize{
		P8:  appleDeveloper.P8,
		Iss: appleDeveloper.Iss,
		Kid: appleDeveloper.Kid,
	}
	authorize.DeleteCertificatesByCerId(appleDeveloper.CerID)
	return nil
}

func (a *AppleDeveloper) List(content string, page, pageSize *int) ([]model.AppleDeveloper, int64, error) {
	appleDevelopers, total, err := a.store.AppleDeveloper.List(content, page, pageSize)
	if err != nil {
		return nil, 0, e.NewWithStack(e.DBError, err)
	}
	return appleDevelopers, total, nil
}

func (a *AppleDeveloper) Update(iss string, limit int, enable bool) error {
	err := a.store.AppleDeveloper.UpdateSetup(iss, limit, enable)
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	return nil
}
