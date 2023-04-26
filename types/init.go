package types

type EncryptPayload struct {
	Ts           int64    `json:"ts"`
	MacAddresses []string `json:"mac_addresses"` //since one device can have more than one network interface (ethernet, wifi ...)
}

type Response struct {
	Encrypted string `json:"mm_encrypted_string"`
	Version   string `json:"mm_encrypted_version"`
}
