package cryptoutil

import "crypto/tls"

// CipherSuites 获取安全的 tls 加密套件
func CipherSuites() (suites []uint16) {
	arr := tls.CipherSuites()
	for _, item := range arr {
		if item.Insecure {
			continue
		}
		suites = append(suites, item.ID)
	}
	return suites
}

// InsecureCipherSuites 获取不安全的 tls 加密套件
func InsecureCipherSuites() (suites []uint16) {
	arr := tls.InsecureCipherSuites()
	for _, item := range arr {
		if !item.Insecure {
			continue
		}
		suites = append(suites, item.ID)
	}
	return suites
}
