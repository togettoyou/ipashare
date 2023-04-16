package sign

import (
	"bytes"
	"ipashare/pkg/tools"
	"os/exec"
)

const (
	serverCrt    = "data/server.crt"
	serverKey    = "data/server.key"
	certChainCrt = "data/cert-chain.crt"
)

func MobileConfig(input []byte) ([]byte, error) {
	if !tools.PathIsExist(serverCrt) ||
		!tools.PathIsExist(serverKey) ||
		!tools.PathIsExist(certChainCrt) {
		return input, nil
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
		return nil, err
	}
	return output.Bytes(), nil
}
