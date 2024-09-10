package dns

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	ib "github.com/infobloxopen/infoblox-go-client/v2"
	"github.com/rs/zerolog/log"
)

func (i *InfobloxProvider) CreateRecord(name, ipAddress string) error {
	objMgr := ib.NewObjectManager(i.connector, "", "")

	result, err := objMgr.CreateHostRecord(
		true,      //enabledns
		false,     //enabledhcp
		name,      //recordName
		DNS_VIEW,  //DNS View
		DNS_ZONE,  // DNS Zone
		"",        //ipv4cidr
		"",        //ipv6cidr
		ipAddress, //ipv4Addr
		"",        //ipv6Addr
		"",        //macAddr
		"",
		true, //useTTL
		30,
		"",
		nil,
		[]string{},
	)

	if err != nil {
		slog.Error("Error creating record", err.Error(), nil)
		return err
	}
	// _, err := i.connector.CreateObject(&record)
	// if err != nil {
	// 	fmt.Printf("Error creating DNS record: %v\n", err)
	// 	return err
	// }
	slog.Info("DNS record created with ", *result.Name, result.DnsName)
	return nil
}

func (i *InfobloxProvider) CreateRecordHttp(name, ipAddress string) error {
	res, err := httpRequest(name, ipAddress, "POST")
	if err != nil {
		slog.Info(err.Error())
		return err
	}

	if res.Body != nil {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Error().Msgf("Error reading response body: %v", err)
			return err
		}
		if res.StatusCode <= 199 || res.StatusCode >= 299 {
			slog.Error(string(data))
			return errors.New(res.Status)
		}
		var record []interface{}
		if err = json.Unmarshal(data, &record); err != nil {
			slog.Error("Error unmarshaling response: ", err.Error(), "")
			return err
		}
	}
	return nil
}

type Record struct {
	Name      string    `json:"name"`
	Ipv4addrs []IpAddrs `json:"ipv4addrs"`
	Comment   string    `json:"comment"`
	View      string    `json:"view"`
}

type IpAddrs struct {
	Ipv4addr string `json:"ipv4addr"`
}

func httpRequest(name, ipAddress, method string) (*http.Response, error) {
	ibs := InfoBloxServer{
		Server:   DNS_SERVER,
		Protocol: "https",
		UserName: DNS_SERVER_USERNAME,
		Password: DNS_SERVER_PASSWORD,
		View:     DNS_VIEW, //Net View
		Zone:     DNS_ZONE, //Zone
	}

	murl := fmt.Sprintf("%s://%s/wapi/v2.12.1/record:%s",
		ibs.Protocol,
		ibs.Server,
		"host")

	slog.Error("Create URL: ", murl, "")
	/*
		curl -k "https://ibtest.vgregion.se/wapi/v2.12.1/record:host" -d '                        ✔
		{"name":"lookie.k8s.vgregion.se",
		"ipv4addrs":[{"ipv4addr":"192.168.10.236"}],
		"comment":"testing dns controller",
		"view":"default.container-dev"}' -u container_dev:hipsterfrilla \
		-H "Content-Type: application/json"
	*/
	name = name + "." + ibs.Zone
	ip := IpAddrs{Ipv4addr: ipAddress}
	record := Record{Name: name,
		Comment: "yaya", View: ibs.View,
		Ipv4addrs: []IpAddrs{ip},
	}
	slog.Info("Record to be sent:", fmt.Sprintf("%v", record), nil)

	recordBytes, err := json.Marshal(record)
	if err != nil {
		slog.Error("Marshalling not mashled ", err.Error(), nil)
		return nil, err
	}
	body := bytes.NewReader(recordBytes)

	req, err := http.NewRequest(method, murl, body)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		slog.Error("Error: something went wrong", err.Error(), nil)
		return nil, err
	}
	//req.SetBasicAuth(os.Getenv("INFOBLOX_USERNAME"), os.Getenv("INFOBLOX_PASSWORD"))
	req.SetBasicAuth(ibs.UserName, ibs.Password)
	req.Header.Set("Content-Type", "application/json")

	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		slog.Error(err.Error(), "", nil)
		return nil, err
	}
	//defer res.Body.Close()
	return res, err
}
