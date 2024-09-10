package dns

import (
	"testing"

	ib "github.com/infobloxopen/infoblox-go-client/v2"
)

func TestInfobloxProviderDeleteRecord(t *testing.T) {
	type fields struct {
		connector *ib.Connector
	}
	type args struct {
		name  string
		value string
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
			if err := i.DeleteRecord(tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("InfobloxProvider.DeleteRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
