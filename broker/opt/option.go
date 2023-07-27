package opt

import (
	"math/rand"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqType uint8

const (
	MQ_TYPE_NONE MqType = iota
	MQ_TYPE_TCP
	MQ_TYPE_WS
)

type Option struct {
	ClientName string
	Name       string
	Pwd        string
	Cert       string
	Key        string
	CaCert     string
}

func (opt *Option) TypeOf(addr string) MqType {
	var addrs []string = strings.Split(addr, ":")
	switch strings.ToLower(addrs[0]) {
	case "tcp":
		return MQ_TYPE_TCP
	case "ws":
		return MQ_TYPE_WS
	default:
		return MQ_TYPE_NONE
	}

}

func (opt *Option) get(addr string) *mqtt.ClientOptions {
	var Opt *mqtt.ClientOptions
	switch opt.TypeOf(addr) {
	case MQ_TYPE_TCP:
		Opt = opt.Credential()
	// case MQ_TYPE_WS:
	// Opt = opt.WebSocket()
	default:
		Opt = opt.None()
	}
	/*
		if opt.CaCert != "" && opt.Cert != "" && opt.Key != "" {
			tlsConf, err := opt.tls()
			if err == nil {
				Opt.SetTLSConfig(tlsConf)
			}
		}
	*/
	return Opt.AddBroker(addr)
}

func (opt *Option) None() *mqtt.ClientOptions {
	if strings.TrimSpace(opt.ClientName) == "" {
		opt.ClientName = genID()
	}

	return &mqtt.ClientOptions{
		ClientID: opt.ClientName,
	}
}

func (opt *Option) Credential() *mqtt.ClientOptions {
	if strings.TrimSpace(opt.ClientName) == "" {
		opt.ClientName = genID()
	}

	return &mqtt.ClientOptions{
		ClientID: opt.ClientName,
		Username: opt.Name,
		Password: opt.Pwd,
	}
}

/*
	공장에서 사용되기에 websocket으로 연결될 경우가 없어 주석처리
	func (opt *Option) WebSocket() *mqtt.ClientOptions {
		return &mqtt.ClientOptions{
			ClientID:         genID(),
			WebsocketOptions: opt.websocket(),
		}
	}

	func (opt *Option) websocket() *mqtt.WebsocketOptions {
		return &mqtt.WebsocketOptions{
			ReadBufferSize: 0,
			WriteBufferSize: 0,
			Proxy: mqtt.ProxyFunction,
		}
	}


// Import trusted certificates from CAfile.pem.
// Alternatively, manually add CA certificates to default openssl CA bundle.
func (opt *Option) tls() (*tls.Config, error) {
	certPool := x509.NewCertPool()
	bytes, err := ioutil.ReadFile(opt.CaCert)
	if err != nil {
		return nil, err
	}

	ok := certPool.AppendCertsFromPEM(bytes)
	if !ok {
		return nil, err
	}

	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair(opt.Cert, opt.Key)
	if err != nil {
		return nil, err
	}

	// Just to print out the client certificate
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certPool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}, nil
}
*/

func genID() string {
	var (
		sLetterRunes []rune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		sSeed        int64  = time.Now().UnixNano()
		sIDLen       int64  = 10
		sRandomID    []rune = make([]rune, sIDLen)
	)

	rand.Seed(sSeed) // random의 seed값을 설정

	for sIndex := range sRandomID {
		sRandomID[sIndex] = sLetterRunes[rand.Intn(len(sLetterRunes))]
	}

	return string(sRandomID)
}
