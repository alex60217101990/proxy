package enums

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type ProtocolType int8

const (
	TCP ProtocolType = iota + 1
	UDP
)

var (
	_ProtocolTypeNameToValue = map[string]ProtocolType{
		"TCP": TCP,
		"tcp": TCP,
		"UDP": UDP,
		"udp": UDP,
	}

	_ProtocolTypeValueToName = map[ProtocolType]string{
		TCP: "tcp",
		UDP: "udp",
	}
)

func (r ProtocolType) MarshalYAML() (interface{}, error) {
	s, ok := _ProtocolTypeValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid ProtocolType: %d", r)
	}
	return s, nil
}

func (r *ProtocolType) UnmarshalYAML(value *yaml.Node) error {
	v, ok := _ProtocolTypeNameToValue[value.Value]
	if !ok {
		return fmt.Errorf("invalid ProtocolType %q", value.Value)
	}
	*r = v
	return nil
}

func (r ProtocolType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _ProtocolTypeValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid ProtocolType: %d", r)
	}
	return json.Marshal(s)
}

func (r *ProtocolType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ProtocolType should be a string, got %s", data)
	}
	v, ok := _ProtocolTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid ProtocolType %q", s)
	}
	*r = v
	return nil
}

func (r ProtocolType) Val() int {
	return int(r)
}

func (r ProtocolType) String() string {
	return _ProtocolTypeValueToName[r]
}
