package be

import (
	"github.com/eltiocaballoloco/sinaloa-cli/helpers"
)

// Implement fucntion to renew haproxy SSL certificates
func RenewHaproxySSLCertificatesStorj(storjPathCert string, storjPathKey string, secret string) []byte {
	// Implement renew haproxy SSL certificates
	return helpers.HandleResponse("Renew haproxy SSL certificates", "200", struct{}{})
}
