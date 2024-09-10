package dns

import (
	"testing"

	ib "github.com/infobloxopen/infoblox-go-client/v2"
)

func TestInfobloxProviderCreateRecord(t *testing.T) {
	type fields struct {
		connector *ib.Connector
	}
	type args struct {
		name      string
		ipAddress string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "CreateRecord",
			args: args{name: "exampleapp", ipAddress: "192.168.1.236"},
		},
		{
			name: "CreateRecord",
			args: args{name: "loggadev", ipAddress: "127.0.0.1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InfobloxProvider{
				connector: tt.fields.connector,
			}
			if err := i.CreateRecord(tt.args.name, tt.args.ipAddress); (err != nil) != tt.wantErr {
				t.Errorf("InfobloxProvider.CreateRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInfobloxProviderCreateRecordHttp(t *testing.T) {
	type fields struct {
		connector *ib.Connector
	}
	type args struct {
		name      string
		ipAddress string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InfobloxProvider{
				connector: tt.fields.connector,
			}
			if err := i.CreateRecordHttp(tt.args.name, tt.args.ipAddress); (err != nil) != tt.wantErr {
				t.Errorf("InfobloxProvider.CreateRecordHttp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
