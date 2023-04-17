package sign

import (
	"bytes"
	"ipashare/pkg/tools"
	"log"
	"os"
	"os/exec"
	"sync"
)

const (
	serverCrt    = "data/server.crt"
	serverKey    = "data/server.key"
	certChainCrt = "data/cert-chain.crt"
)

var crtAndKeyInfoMux sync.RWMutex

type CrtAndKeyInfo struct {
	ServerCrtContent    string `json:"server_crt_content"`
	ServerKeyContent    string `json:"server_key_content"`
	CertChainCrtContent string `json:"cert_chain_crt_content"`
}

func SetCrtAndKey(crtAndKeyInfo *CrtAndKeyInfo) error {
	crtAndKeyInfoMux.Lock()
	defer crtAndKeyInfoMux.Unlock()

	err := write(serverCrt, []byte(crtAndKeyInfo.ServerCrtContent))
	if err != nil {
		return err
	}
	err = write(serverKey, []byte(crtAndKeyInfo.ServerKeyContent))
	if err != nil {
		return err
	}
	return write(certChainCrt, []byte(crtAndKeyInfo.CertChainCrtContent))
}

func write(filename string, data []byte) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func GetCrtAndKey() CrtAndKeyInfo {
	crtAndKeyInfoMux.RLock()
	defer crtAndKeyInfoMux.RUnlock()

	if !tools.PathIsExist(serverCrt) ||
		!tools.PathIsExist(serverKey) ||
		!tools.PathIsExist(certChainCrt) {
		return CrtAndKeyInfo{}
	}

	serverCrtContent, _ := os.ReadFile(serverCrt)
	serverKeyContent, _ := os.ReadFile(serverKey)
	certChainCrtContent, _ := os.ReadFile(certChainCrt)
	return CrtAndKeyInfo{
		ServerCrtContent:    string(serverCrtContent),
		ServerKeyContent:    string(serverKeyContent),
		CertChainCrtContent: string(certChainCrtContent),
	}
}

func MobileConfig(input []byte) []byte {
	crtAndKey := GetCrtAndKey()
	if crtAndKey.ServerCrtContent == "" ||
		crtAndKey.ServerKeyContent == "" ||
		crtAndKey.CertChainCrtContent == "" {
		return input
	}

	var output bytes.Buffer
	opensslCmd := exec.Command("openssl", "smime", "-sign",
		"-in", "/dev/stdin",
		"-out", "/dev/stdout",
		"-signer", serverCrt,
		"-inkey", serverKey,
		"-certfile", certChainCrt,
		"-outform", "der", "-nodetach")
	opensslCmd.Stdin = bytes.NewBuffer(input)
	opensslCmd.Stdout = &output

	err := opensslCmd.Run()
	if err != nil {
		log.Println("mobileconfig签名异常：", err)
		return input
	}
	return output.Bytes()
}
