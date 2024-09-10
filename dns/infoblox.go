package dns

import (
	"log/slog"
	"os"

	ib "github.com/infobloxopen/infoblox-go-client/v2"
)

type InfobloxProvider struct {
	connector *ib.Connector
}

type Infoblox struct {
	InfobloxServer        string
	InfobloxServerPort    string
	InfobloxUsername      string
	InfobloxPassword      string
	InfoBloxServerVersion string
}

var (
	DNS_SERVER          = os.Getenv("DNS_SERVER")
	DNS_SERVER_PORT     = os.Getenv("DNS_SERVER_PORT")
	DNS_SERVER_USERNAME = os.Getenv("DNS_SERVER_USERNAME")
	DNS_SERVER_PASSWORD = os.Getenv("DNS_SERVER_PASSWORD")
	DNS_VERSION         = os.Getenv("DNS_VERSION")
	DNS_VIEW            = os.Getenv("DNS_VIEW") //default.container-dev
	DNS_ZONE            = os.Getenv("DNS_ZONE") //k8s.localdev.me
)

func NewInfobloxProvider() (*InfobloxProvider, error) {
	hc := getHostConfig()
	tc := getTransportConfig()
	wrb := &ib.WapiRequestBuilder{}
	whr := &ib.WapiHttpRequestor{}
	ac := getAuthConfig()
	client, err := ib.NewConnector(hc, ac, tc, wrb, whr)
	if err != nil {
		slog.Error("Error creating a connector for infoblox", err.Error(), nil)
		return nil, err
	}

	return &InfobloxProvider{
		connector: client,
	}, nil
}

func NewInfoBloxServer(server, port, version string) Infoblox {
	return Infoblox{
		InfobloxServer:        server,
		InfobloxServerPort:    port,
		InfoBloxServerVersion: version,
	}
}

func getHostConfig() ib.HostConfig {
	return ib.HostConfig{
		Scheme:  "https",
		Host:    DNS_SERVER,
		Version: DNS_VERSION,
		Port:    DNS_SERVER_PORT,
	}
}

func getTransportConfig() ib.TransportConfig {
	// tc := ib.TransportConfig{
	// 	SslVerify: false,
	// 	HttpRequestTimeout: 20,
	// 	HttpPoolConnections: 10,
	// }
	return ib.NewTransportConfig("false", 20, 10)
}

func getAuthConfig() ib.AuthConfig {
	return ib.AuthConfig{
		Username: DNS_SERVER_USERNAME,
		Password: DNS_SERVER_PASSWORD,
	}
}
