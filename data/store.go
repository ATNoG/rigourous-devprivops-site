package data

import (
	"encoding/json"
	"os"
)

type Store struct {
	Data []*Report
}

func FromFile(file string) (*Store, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	payload := new(Store)
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

/*
func (ds *Store) SetAccess(level int, group string) {
	ds.Level = level
	ds.Group = group
}
*/

func (s *Store) ToFile(file string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(file, data, 0666); err != nil {
		return err
	}

	return nil
}
