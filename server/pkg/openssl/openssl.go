package openssl

import (
	"fmt"
	"ipashare/pkg/tools"
)

const (
	bash                 = "/bin/bash"
	keyAndReqCSRConfPath = "KeyAndReqCSR.conf"
)

func GenKeyAndReqCSR(keyPath, csrPath string) error {
	err := tools.CmdClient.Command(bash, "-c",
		fmt.Sprintf("openssl genrsa -out %s 2048", keyPath),
	)
	if err != nil {
		return err
	}
	if !tools.PathIsExist(keyAndReqCSRConfPath) {
		err = tools.CreateFile(`[ req ]
prompt = no
distinguished_name = req_distinguished

[ req_distinguished ]
C = CN
O = ipaShare Company
CN = ipaShare`, keyAndReqCSRConfPath)
		if err != nil {
			return err
		}
	}
	return tools.CmdClient.Command(bash, "-c",
		fmt.Sprintf("openssl req -new -config %s -sha256 -key %s -out %s",
			keyAndReqCSRConfPath, keyPath, csrPath),
	)
}

func GenPEM(cerPath, pemPath string) error {
	return tools.CmdClient.Command(bash, "-c",
		fmt.Sprintf("openssl x509 -in %s -inform DER -outform PEM -out %s",
			cerPath, pemPath),
	)
}
