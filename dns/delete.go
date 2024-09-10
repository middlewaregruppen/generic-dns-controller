package dns

import "log/slog"

func (i *InfobloxProvider) DeleteRecord(name, value string) error {
	_, err := i.connector.DeleteObject(name)
	if err != nil {
		slog.Error("Error deleting DNS record", err.Error(), nil)
		return err
	}
	return nil
}
