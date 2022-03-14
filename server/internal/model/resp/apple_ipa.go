package resp

import "ipashare/internal/model"

type AppleIPA struct {
	model.AppleIPA
	IconUrl    string `json:"icon_url"`
	InstallUrl string `json:"install_url"`
}
