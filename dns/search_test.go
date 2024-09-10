package dns

import (
	"reflect"
	"testing"

	ib "github.com/infobloxopen/infoblox-go-client/v2"
)

func TestInfobloxProviderSearchRecordHttp(t *testing.T) {
	type fields struct {
		connector *ib.Connector
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Record
		wantErr bool
	}{
		{
			name:    "SearchRecordHttp",
			args:    args{name: "exampleapp"},
			want:    Record{Name: "exampleapp", Ipv4addrs: []IpAddrs{{"192.168.1.236"}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InfobloxProvider{
				connector: tt.fields.connector,
			}
			got, err := i.SearchRecordHttp(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("InfobloxProvider.SearchRecordHttp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InfobloxProvider.SearchRecordHttp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfobloxProviderSearchRecord(t *testing.T) {
	type fields struct {
		connector *ib.Connector
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "SearchRecord",
			args:    args{name: "exampleapp"},
			want:    true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InfobloxProvider{
				connector: tt.fields.connector,
			}
			got, err := i.SearchRecord(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("InfobloxProvider.SearchRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InfobloxProvider.SearchRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}
