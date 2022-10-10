package appstore

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

type Authorize struct {
	P8  string
	Iss string
	Kid string
}

type Links struct {
	Self string `json:"self"`
}

type Meta struct {
	Paging struct {
		Total int `json:"total"`
		Limit int `json:"limit"`
	} `json:"paging"`
}

type DevicesData struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		AddedDate   string `json:"addedDate"`
		Name        string `json:"name"`
		DeviceClass string `json:"deviceClass"`
		Model       string `json:"model"`
		Udid        string `json:"udid"`
		Platform    string `json:"platform"`
		Status      string `json:"status"`
	} `json:"attributes"`
	Links Links `json:"links"`
}

type DevicesResponse struct {
	Data  DevicesData `json:"data"`
	Links Links       `json:"links"`
}

type DevicesResponseList struct {
	Data  []DevicesData `json:"data"`
	Links Links         `json:"links"`
	Meta  Meta          `json:"meta"`
}

type CertificateData struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		SerialNumber       string      `json:"serialNumber"`
		CertificateContent string      `json:"certificateContent"`
		DisplayName        string      `json:"displayName"`
		Name               string      `json:"name"`
		CsrContent         interface{} `json:"csrContent"`
		Platform           string      `json:"platform"`
		ExpirationDate     string      `json:"expirationDate"`
		CertificateType    string      `json:"certificateType"`
	} `json:"attributes"`
	Links Links `json:"links"`
}

type CertificateResponse struct {
	Data  CertificateData `json:"data"`
	Links Links           `json:"links"`
}

type CertificateResponseList struct {
	Data  []CertificateData `json:"data"`
	Links Links             `json:"links"`
	Meta  Meta              `json:"meta"`
}

type BundleIdData struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		Name       string `json:"name"`
		Identifier string `json:"identifier"`
		Platform   string `json:"platform"`
		SeedID     string `json:"seedId"`
	} `json:"attributes"`
	Relationships struct {
		BundleIDCapabilities struct {
			Meta struct {
				Paging struct {
					Total int   `json:"total"`
					Limit int64 `json:"limit"`
				} `json:"paging"`
			} `json:"meta"`
			Data  []interface{} `json:"data"`
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"bundleIdCapabilities"`
		Profiles struct {
			Meta struct {
				Paging struct {
					Total int   `json:"total"`
					Limit int64 `json:"limit"`
				} `json:"paging"`
			} `json:"meta"`
			Data  []interface{} `json:"data"`
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"profiles"`
	} `json:"relationships"`
	Links Links `json:"links"`
}

type BundleIdResponse struct {
	Data  BundleIdData `json:"data"`
	Links Links        `json:"links"`
}

type BundleIdResponseList struct {
	Data  []BundleIdData `json:"data"`
	Links Links          `json:"links"`
	Meta  Meta           `json:"meta"`
}

type ProfileData struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		ProfileState   string `json:"profileState"`
		CreatedDate    string `json:"createdDate"`
		ProfileType    string `json:"profileType"`
		Name           string `json:"name"`
		ProfileContent string `json:"profileContent"`
		UUID           string `json:"uuid"`
		Platform       string `json:"platform"`
		ExpirationDate string `json:"expirationDate"`
	} `json:"attributes"`
	Relationships struct {
		BundleID struct {
			Data struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"data"`
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"bundleId"`
		Certificates struct {
			Meta struct {
				Paging struct {
					Total int   `json:"total"`
					Limit int64 `json:"limit"`
				} `json:"paging"`
			} `json:"meta"`
			Data []struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"data"`
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"certificates"`
		Devices struct {
			Meta struct {
				Paging struct {
					Total int   `json:"total"`
					Limit int64 `json:"limit"`
				} `json:"paging"`
			} `json:"meta"`
			Data []struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"data"`
			Links struct {
				Self    string `json:"self"`
				Related string `json:"related"`
			} `json:"links"`
		} `json:"devices"`
	} `json:"relationships"`
	Links Links `json:"links"`
}

type ProfileResponse struct {
	Data  ProfileData `json:"data"`
	Links Links       `json:"links"`
}

const (
	devicesUrl      = "https://api.appstoreconnect.apple.com/v1/devices"
	certificatesUrl = "https://api.appstoreconnect.apple.com/v1/certificates"
	bundleIdsUrl    = "https://api.appstoreconnect.apple.com/v1/bundleIds"
	profilesUrl     = "https://api.appstoreconnect.apple.com/v1/profiles"
)

// GetAvailableDevices 获取账号可用测试设备列表
func (auth *Authorize) GetAvailableDevices() (*DevicesResponseList, error) {
	// TODO: 暂时通过多次请求解决403问题
	for i := 0; i < 15; i++ {
		resp, err := auth.httpRequest("GET", devicesUrl+"?limit=200", nil)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode() == 200 {
			var devicesResponseList DevicesResponseList
			err = json.Unmarshal(resp.Body(), &devicesResponseList)
			if err != nil {
				return nil, err
			}
			fasthttp.ReleaseResponse(resp)
			return &devicesResponseList, nil
		}
		if resp.StatusCode() != 403 {
			fasthttp.ReleaseResponse(resp)
			return nil, errors.New(fmt.Sprintf("%s", resp.Body()))
		}
		fasthttp.ReleaseResponse(resp)
	}
	// 15次请求还是403后 再来一次还是403就直接报错
	resp, err := auth.httpRequest("GET", devicesUrl+"?limit=200", nil)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New(fmt.Sprintf("%s", resp.Body()))
	}
	var devicesResponseList DevicesResponseList
	err = json.Unmarshal(resp.Body(), &devicesResponseList)
	if err != nil {
		return nil, err
	}
	return &devicesResponseList, nil
}

// GetAvailableDevice 通过 deviceID 获取账号指定测试设备
func (auth *Authorize) GetAvailableDevice(deviceID string) (*DevicesResponse, bool, error) {
	resp, err := auth.httpRequest("GET", devicesUrl+"/"+deviceID, nil)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return nil, false, err
	}
	if resp.StatusCode() == 404 {
		return nil, false, nil
	}
	if resp.StatusCode() == 200 {
		var devicesResponse DevicesResponse
		err = json.Unmarshal(resp.Body(), &devicesResponse)
		if err != nil {
			return nil, false, err
		}
		return &devicesResponse, true, nil
	}
	return nil, false, errors.New(fmt.Sprintf("%s", resp.Body()))
}

// GetAvailableDeviceByUDID 通过 udid 获取账号指定测试设备
func (auth *Authorize) GetAvailableDeviceByUDID(udid string) (*DevicesResponse, bool, error) {
	resp, err := auth.httpRequest("GET", devicesUrl+"?limit=200&filter[udid]="+udid, nil)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return nil, false, err
	}
	if resp.StatusCode() != 200 {
		return nil, false, errors.New(fmt.Sprintf("%s", resp.Body()))
	}
	var devicesResponseList DevicesResponseList
	err = json.Unmarshal(resp.Body(), &devicesResponseList)
	if err != nil {
		return nil, false, err
	}
	if len(devicesResponseList.Data) != 1 {
		return nil, false, nil
	}
	return &DevicesResponse{
		Data:  devicesResponseList.Data[0],
		Links: devicesResponseList.Links,
	}, true, nil
}

// AddAvailableDevice 添加测试设备
func (auth *Authorize) AddAvailableDevice(udid string) (*DevicesResponse, error) {
	type DeviceCreateRequest struct {
		Data struct {
			Type       string `json:"type"`
			Attributes struct {
				Name     string `json:"name"`
				Udid     string `json:"udid"`
				Platform string `json:"platform"`
			} `json:"attributes"`
		} `json:"data"`
	}
	var deviceCreateRequest DeviceCreateRequest
	deviceCreateRequest.Data.Type = "devices"
	deviceCreateRequest.Data.Attributes.Name = udid
	deviceCreateRequest.Data.Attributes.Udid = udid
	deviceCreateRequest.Data.Attributes.Platform = "IOS"
	jsonStr, err := json.Marshal(&deviceCreateRequest)
	if err != nil {
		return nil, err
	}
	resp, err := auth.httpRequest("POST", devicesUrl, jsonStr)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 201 {
		return nil, errors.New(fmt.Sprintf("%s", resp.Body()))
	}
	var devicesResponse DevicesResponse
	err = json.Unmarshal(resp.Body(), &devicesResponse)
	if err != nil {
		return nil, err
	}
	return &devicesResponse, nil
}

// CreateCertificates 创建证书
func (auth *Authorize) CreateCertificates(csr string) (*CertificateResponse, error) {
	type CertificateCreateRequest struct {
		Data struct {
			Type       string `json:"type"`
			Attributes struct {
				CsrContent      string `json:"csrContent"`
				CertificateType string `json:"certificateType"`
			} `json:"attributes"`
		} `json:"data"`
	}
	var certificateCreateRequest CertificateCreateRequest
	certificateCreateRequest.Data.Type = "certificates"
	certificateCreateRequest.Data.Attributes.CertificateType = "IOS_DEVELOPMENT"
	certificateCreateRequest.Data.Attributes.CsrContent = csr
	jsonStr, err := json.Marshal(&certificateCreateRequest)
	if err != nil {
		return nil, err
	}
	resp, err := auth.httpRequest("POST", certificatesUrl, jsonStr)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 201 {
		return nil, errors.New(fmt.Sprintf("%s", resp.Body()))
	}
	var certificateResponse CertificateResponse
	err = json.Unmarshal(resp.Body(), &certificateResponse)
	if err != nil {
		return nil, err
	}
	return &certificateResponse, nil
}

// GetCertificatesList 获取账号下所有证书
func (auth *Authorize) GetCertificatesList() ([]CertificateData, error) {
	resp, err := auth.httpRequest("GET", certificatesUrl, nil)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New(fmt.Sprintf("%s", resp.Body()))
	}
	var certificateResponseList CertificateResponseList
	err = json.Unmarshal(resp.Body(), &certificateResponseList)
	if err != nil {
		return nil, err
	}

	certificateDataList := make([]CertificateData, 0)
	for _, v := range certificateResponseList.Data {
		if v.Type == "certificates" &&
			(v.Attributes.CertificateType == "IOS_DEVELOPMENT" ||
				v.Attributes.CertificateType == "MAC_APP_DEVELOPMENT" ||
				v.Attributes.CertificateType == "DEVELOPMENT") {
			certificateDataList = append(certificateDataList, v)
		}
	}
	return certificateDataList, nil
}

// DeleteAllCertificates 删除账号下所有证书
func (auth *Authorize) DeleteAllCertificates() error {
	resp, err := auth.httpRequest("GET", certificatesUrl, nil)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("%s", resp.Body()))
	}
	var certificateResponseList CertificateResponseList
	err = json.Unmarshal(resp.Body(), &certificateResponseList)
	if err != nil {
		return err
	}

	for _, v := range certificateResponseList.Data {
		resp, err := auth.httpRequest("DELETE", certificatesUrl+"/"+v.ID, nil)
		if err != nil {
			return err
		}
		if resp.StatusCode() != 204 {
			return errors.New(fmt.Sprintf("DeleteAllCertificates %s", resp.Body()))
		}
		fasthttp.ReleaseResponse(resp)
	}
	return nil
}

// DeleteCertificatesByCerId 删除账号下指定的证书
func (auth *Authorize) DeleteCertificatesByCerId(cerId string) error {
	resp, err := auth.httpRequest("DELETE", certificatesUrl+"/"+cerId, nil)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 204 {
		return errors.New(fmt.Sprintf("DeleteCertificatesByCerId %s", resp.Body()))
	}
	return nil
}

// CreateBundleIds 创建BundleIds
func (auth *Authorize) CreateBundleIds(identifier string) (*BundleIdResponse, error) {
	type BundleIdCreateRequest struct {
		Data struct {
			Type       string `json:"type"`
			Attributes struct {
				Name       string `json:"name"`
				Platform   string `json:"platform"`
				SeedID     string `json:"seedId"`
				Identifier string `json:"identifier"`
			} `json:"attributes"`
		} `json:"data"`
	}
	var bundleIdCreateRequest BundleIdCreateRequest
	bundleIdCreateRequest.Data.Type = "bundleIds"
	bundleIdCreateRequest.Data.Attributes.Name = "AppBundleId"
	bundleIdCreateRequest.Data.Attributes.Identifier = identifier
	bundleIdCreateRequest.Data.Attributes.Platform = "IOS"
	jsonStr, err := json.Marshal(&bundleIdCreateRequest)
	if err != nil {
		return nil, err
	}
	resp, err := auth.httpRequest("POST", bundleIdsUrl, jsonStr)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 201 {
		return nil, errors.New(fmt.Sprintf("%s", resp.Body()))
	}
	var bundleIdResponse BundleIdResponse
	err = json.Unmarshal(resp.Body(), &bundleIdResponse)
	if err != nil {
		return nil, err
	}
	return &bundleIdResponse, nil
}

// GetBundleIdsByIdentifier 根据identifier获取BundleIds
func (auth *Authorize) GetBundleIdsByIdentifier(identifier string) (string, error) {
	resp, err := auth.httpRequest("GET", bundleIdsUrl, nil)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", errors.New(fmt.Sprintf("%s", resp.Body()))
	}
	var bundleIdResponseList BundleIdResponseList
	err = json.Unmarshal(resp.Body(), &bundleIdResponseList)
	if err != nil {
		return "", err
	}
	for _, v := range bundleIdResponseList.Data {
		if v.Attributes.Identifier == identifier {
			return v.ID, nil
		}
	}
	return "", nil
}

// DeleteAllBundleIds 删除账号下所有的BundleIds
func (auth *Authorize) DeleteAllBundleIds() error {
	resp, err := auth.httpRequest("GET", bundleIdsUrl, nil)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("%s", resp.Body()))
	}
	var bundleIdResponseList BundleIdResponseList
	err = json.Unmarshal(resp.Body(), &bundleIdResponseList)
	if err != nil {
		return err
	}

	for _, v := range bundleIdResponseList.Data {
		resp, err := auth.httpRequest("DELETE", bundleIdsUrl+"/"+v.ID, nil)
		if err != nil {
			return err
		}
		if resp.StatusCode() != 204 {
			return errors.New(fmt.Sprintf("DeleteAllBundleIds %s", resp.Body()))
		}
		fasthttp.ReleaseResponse(resp)
	}
	return nil
}

// CreateProfile 创建Profile
func (auth *Authorize) CreateProfile(name string, bundleID string, cerId string, devicesId string) (*ProfileResponse, error) {
	type Data struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}
	type ProfileCreateRequest struct {
		Data struct {
			Type       string `json:"type"`
			Attributes struct {
				Name        string `json:"name"`
				ProfileType string `json:"profileType"`
			} `json:"attributes"`
			Relationships struct {
				BundleID struct {
					Data Data `json:"data"`
				} `json:"bundleId"`
				Certificates struct {
					Data []Data `json:"data"`
				} `json:"certificates"`
				Devices struct {
					Data []Data `json:"data"`
				} `json:"devices"`
			} `json:"relationships"`
		} `json:"data"`
	}
	var profileCreateRequest ProfileCreateRequest
	profileCreateRequest.Data.Type = "profiles"
	profileCreateRequest.Data.Attributes.ProfileType = "IOS_APP_DEVELOPMENT"
	profileCreateRequest.Data.Attributes.Name = name
	profileCreateRequest.Data.Relationships.BundleID.Data.Type = "bundleIds"
	profileCreateRequest.Data.Relationships.BundleID.Data.ID = bundleID
	profileCreateRequest.Data.Relationships.Certificates.Data = []Data{
		{
			ID:   cerId,
			Type: "certificates",
		},
	}
	profileCreateRequest.Data.Relationships.Devices.Data = []Data{
		{
			ID:   devicesId,
			Type: "devices",
		},
	}
	jsonStr, err := json.Marshal(&profileCreateRequest)
	if err != nil {
		return nil, err
	}
	resp, err := auth.httpRequest("POST", profilesUrl, jsonStr)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 201 {
		return nil, errors.New(fmt.Sprintf("%s", resp.Body()))
	}
	var profileResponse ProfileResponse
	err = json.Unmarshal(resp.Body(), &profileResponse)
	if err != nil {
		return nil, err
	}
	return &profileResponse, nil
}

func (auth *Authorize) httpRequest(method string, url string, body []byte) (*fasthttp.Response, error) {
	token, err := auth.createToken()
	if err != nil {
		return nil, err
	}
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetContentType("application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.SetMethod(method)
	req.SetRequestURI(url)
	req.SetBody(body)
	if err := fasthttp.Do(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (auth *Authorize) createToken() (string, error) {
	token := &jwt.Token{
		Header: map[string]interface{}{
			"alg": "ES256",
			"kid": auth.Kid,
		},
		Claims: jwt.MapClaims{
			"iss": auth.Iss,
			"exp": time.Now().Add(time.Second * 60 * 5).Unix(),
			"aud": "appstoreconnect-v1",
		},
		Method: jwt.SigningMethodES256,
	}
	privateKey, err := authKeyFromBytes([]byte(auth.P8))
	if err != nil {
		return "", err
	}
	return token.SignedString(privateKey)
}

func authKeyFromBytes(key []byte) (*ecdsa.PrivateKey, error) {
	var err error
	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errors.New("token: AuthKey must be a valid .p8 PEM file")
	}
	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		return nil, err
	}
	var pkey *ecdsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*ecdsa.PrivateKey); !ok {
		return nil, errors.New("token: AuthKey must be of type ecdsa.PrivateKey")
	}
	return pkey, nil
}
