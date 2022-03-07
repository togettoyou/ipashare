package resp

type CertificateData struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		DisplayName     string `json:"display_name"`
		Name            string `json:"name"`
		ExpirationDate  string `json:"expiration_date"`
		CertificateType string `json:"certificate_type"`
	} `json:"attributes"`
}
