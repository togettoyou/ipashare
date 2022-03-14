package sign

import "ipashare/pkg/tools"

func run(pemPath, keyPath, mobileprovisionPath, outputIPAPath, inputIPAPath string) error {
	return tools.CmdClient.Command("zsign",
		"-c", pemPath,
		"-k", keyPath,
		"-m", mobileprovisionPath,
		"-o", outputIPAPath,
		"-z", "9",
		inputIPAPath)
}

func Version() (string, error) {
	output, err := tools.CmdClient.Output("zsign", "-v")
	if err != nil {
		return "", err
	}
	return string(output), nil
}
