package dns

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	ib "github.com/infobloxopen/infoblox-go-client/v2"
)

func (i *InfobloxProvider) UpdateRecord(name string) error {
	fmt.Printf("Searching for Gardot: %s\n", name)
	ref, err := i.SearchRecordHttp(name)
	if err != nil {
		slog.Error("Error searching for record", err.Error(), nil)
		return err
	}
	fmt.Printf("REcOrD: %+v\n", ref)
	v := reflect.ValueOf(ref)
	fmt.Printf("ValueOf: %v\n", v)
	fmt.Printf("Tipe: %v\n", v.Elem())
	iter := reflect.ValueOf(ref).MapRange()
	for iter.Next() {
		if strings.Compare(iter.Key().String(), "ref") == 0 {
			objMgr := ib.NewObjectManager(i.connector, "", "")
			refObj, err := objMgr.UpdateHostRecord(iter.Value().String(),
				true,
				false,
				name,
				DNS_VIEW,
				DNS_ZONE,
				"",
				"",
				"127.0.0.1",
				"",
				"",
				"",
				true,
				30,
				"",
				nil,
				[]string{},
			)
			if err != nil {
				slog.Error("Error updating record", err.Error(), nil)
				return err
			}
			slog.Info("Record updated", refObj.Ref, nil)
			return nil
		}

	}
	return err
}
