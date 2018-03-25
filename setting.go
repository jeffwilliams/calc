package main

import "fmt"

type Setting interface {
	Set(s string) error
}

var Settings = map[string]Setting{}

func init() {
	Settings["obase"] = &outputBase
}

func SetSetting(name, value string) error {
	s, ok := Settings[name]
	if !ok {
		return fmt.Errorf("No sucvh setting %s", name)
	}

	err := s.Set(value)
	if err != nil {
		return err
	}

	updateAutocomplete()
	return nil
}

func SettingExists(name string) (ok bool) {
	_, ok = Settings[name]
	return
}
