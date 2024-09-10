package dns

import (
	"testing"

	ib "github.com/infobloxopen/infoblox-go-client/v2"
)

func TestInfobloxProviderUpdateRecord(t *testing.T) {
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
		wantErr bool
	}{
		{
			name: "UpdateRecord",
			args: args{name: "exampleapp"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InfobloxProvider{
				connector: tt.fields.connector,
			}
			if err := i.UpdateRecord(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("InfobloxProvider.UpdateRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
