package enums

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type OperationType int8

const (
	ADD OperationType = iota + 1
	DELETE
	REBOOT
)

var (
	_OperationTypeNameToValue = map[string]OperationType{
		"ADD":    ADD,
		"add":    ADD,
		"DELETE": DELETE,
		"delete": DELETE,
		"reboot": REBOOT,
		"REBOOT": REBOOT,
	}

	_OperationTypeValueToName = map[OperationType]string{
		ADD:    "add",
		DELETE: "delete",
		REBOOT: "reboot",
	}
)

func (o OperationType) MarshalYAML() (interface{}, error) {
	s, ok := _OperationTypeValueToName[o]
	if !ok {
		return nil, fmt.Errorf("invalid OperationType: %d", o)
	}
	return s, nil
}

func (o *OperationType) UnmarshalYAML(value *yaml.Node) error {
	v, ok := _OperationTypeNameToValue[value.Value]
	if !ok {
		return fmt.Errorf("invalid OperationType %q", value.Value)
	}
	*o = v
	return nil
}

func (o OperationType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(o).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _OperationTypeValueToName[o]
	if !ok {
		return nil, fmt.Errorf("invalid OperationType: %d", o)
	}
	return json.Marshal(s)
}

func (o *OperationType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("OperationType should be a string, got %s", data)
	}
	v, ok := _OperationTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid OperationType %q", s)
	}
	*o = v
	return nil
}

func (o OperationType) Val() int {
	return int(o)
}

func (o OperationType) String() string {
	return _OperationTypeValueToName[o]
}
