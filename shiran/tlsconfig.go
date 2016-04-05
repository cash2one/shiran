package shiran 

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"github.com/golang/glog"
)

func LoadCertificates(certificateFile, privateKeyFile, caFile string) (tls.Certificate, *x509.CertPool) {
	mycert, err := tls.LoadX509KeyPair(certificateFile, privateKeyFile)
	if err != nil {
		glog.Errorf("LoadCertificates LoadX509KeyPair %s %s %s failed %v", certificateFile, privateKeyFile, caFile, err)
		panic(err)	
	}

	pem, err := ioutil.ReadFile(caFile)
	if err != nil {
		glog.Errorf("LoadCertificates ReadFile %s failed %v", caFile, err)
		panic(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pem) {
		glog.Errorf("LoadCertificates appending certs %s %s %s failed %v", certificateFile, privateKeyFile, caFile, err)
		panic(err)
	}

	return mycert, certPool
}

func GetServerTlsConfiguration(certificateFile, privateKeyFile, caFile string) *tls.Config {
	config := &tls.Config{}
	mycert, certPool := LoadCertificates(certificateFile, privateKeyFile, caFile)
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0] = mycert

	config.RootCAs = certPool
	config.ClientCAs = certPool

	//config.ClientAuth = tls.RequireAndVerifyClientCert

	//Optional stuff

	//Use only modern ciphers
	config.CipherSuites = []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256}

	//Use only TLS v1.2
	config.MinVersion = tls.VersionTLS12

	//Don't allow session resumption
	config.SessionTicketsDisabled = true
	return config
}

func GetClientTlsConfiguration(caFile string) *tls.Config {
	config := &tls.Config{}

	pem, err := ioutil.ReadFile(caFile)
	if err != nil {
		glog.Errorf("GetClientTlsConfiguration ReadFile %s failed %v", caFile, err)
		panic(err)	
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pem) {
		glog.Errorf("GetClientTlsConfiguration appending certs %s failed %v", caFile, err)
		panic(err)
	}

	config.RootCAs = certPool
	config.ClientCAs = certPool

	//config.ClientAuth = tls.RequireAndVerifyClientCert

	//Optional stuff

	//Use only modern ciphers
	config.CipherSuites = []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256}

	//Use only TLS v1.2
	config.MinVersion = tls.VersionTLS12

	//Don't allow session resumption
	config.SessionTicketsDisabled = true
	return config
}
