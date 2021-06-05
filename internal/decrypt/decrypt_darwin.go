package decrypt

func ChromePass(key, encryptPass []byte) ([]byte, error) {
	if len(encryptPass) > 3 {
		if len(key) == 0 {
			return nil, errSecurityKeyIsEmpty
		}
		var iv = []byte{32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32}
		return aes128CBCDecrypt(key, iv, encryptPass[3:])
	} else {
		return nil, errDecryptFailed
	}
}

// DPAPI TODO: ReplaceDPAPI
func DPAPI(data []byte) ([]byte, error) {
	return nil, nil
}
