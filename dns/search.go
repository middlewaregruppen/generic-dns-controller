package dns

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	ib "github.com/infobloxopen/infoblox-go-client/v2"
)

func (i *InfobloxProvider) SearchRecord(name string) (bool, error) {
	objMr := ib.NewObjectManager(i.connector, "", "")
	name = name + "." + "k8s-vgregion.se"
	rec, err := objMr.GetHostRecord(DNS_VIEW, DNS_ZONE, name, "192.168.10.236", "")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			slog.Error("Record not found")
			return false, nil
		}
		slog.Error("Search Error", err.Error(), nil)
		return false, err
	}
	slog.Info("Found record for ", rec.DnsName, nil)
	return true, nil
}

type InfoBloxServer struct {
	Server   string `json:"server,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Insecure bool   `json:"insecure,omitempty"`
	UserName string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Zone     string `json:"filter,omitempty"`
	View     string `json:"view,omitempty"`
}

var (
	GRID_URL = "https://" + os.Getenv("DNS_SERVER") + "/wapi/" + os.Getenv("DNS_VERSION") + "/"
)

// Define structs for Record data from search operation
type IPv4Addr struct {
	Ref              string `json:"_ref"`
	IPv4Addr         string `json:"ipv4addr"`
	ConfigureForDHCP bool   `json:"configure_for_dhcp"`
	Host             string `json:"host"`
}

type SeaRecord struct {
	Ref       string     `json:"_ref"`
	IPv4Addrs []IPv4Addr `json:"ipv4addrs"`
	Name      string     `json:"name"`
	View      string     `json:"view"`
}

func (i *InfobloxProvider) SearchRecordHttp(name string) (SeaRecord, error) {
	ibs := InfoBloxServer{
		Server:   DNS_SERVER,
		Protocol: "https",
		UserName: DNS_SERVER_USERNAME,
		Password: DNS_SERVER_PASSWORD,
		View:     DNS_VIEW,
		Zone:     DNS_ZONE,
	}
	var sname string
	if strings.HasSuffix(name, ibs.Zone) {
		sname = name
	} else {
		sname = strings.Split(name, ".")[0] + "." + ibs.Zone
	}

	url := fmt.Sprintf("%s://%s/wapi/v2.12.1/record:%s?name=%s&zone=%s&view=%s",
		ibs.Protocol,
		ibs.Server,
		"host",
		sname,
		ibs.Zone,
		ibs.View)
	slog.Info("Search URL", url, nil)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("Error: something went wrong", err.Error(), nil)
		return SeaRecord{}, err
	}

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
		slog.Error("Error executing request", err.Error(), nil)
		return SeaRecord{}, err
	}

	// Check record does not exist
	if res.StatusCode == 404 {
		slog.Info("Reponse deux", fmt.Sprintf("%d", res.StatusCode), nil)
		return SeaRecord{}, errors.New(res.Status)
	}
	//var record []interface{}
	var searecord SeaRecord

	if res.Body != nil {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			slog.Error("Error reading response body", err.Error(), nil)
			return SeaRecord{}, err
		}

		if res.StatusCode <= 199 || res.StatusCode >= 299 {
			return SeaRecord{}, errors.New(res.Status)
		}
		if err = json.Unmarshal(data, &searecord); err != nil {
			slog.Error("Error unmarshaling response", err.Error(), nil)
			return SeaRecord{}, err
		}
	}
	// if reflect.ValueOf(record).IsNil() {
	// 	return SeaRecord{}, err
	// }

	if len(searecord.Ref) > 0 {
		slog.Info(name, "Record already exists", nil)
		return searecord, nil
	}
	return SeaRecord{}, errors.New("no record found")
}
