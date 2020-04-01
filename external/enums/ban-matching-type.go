package enums

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type BanMatchingType int8

const (
	Strict BanMatchingType = iota + 1
	ByPattern
)

var (
	_BanMatchingTypeNameToValue = map[string]BanMatchingType{
		"Strict":     Strict,
		"ByPattern":  ByPattern,
		"strict":     Strict,
		"by_pattern": ByPattern,
	}

	_BanMatchingTypeValueToName = map[BanMatchingType]string{
		ByPattern: "ByPattern",
		Strict:    "Strict",
	}
)

func (bm BanMatchingType) MarshalYAML() (interface{}, error) {
	s, ok := _BanMatchingTypeValueToName[bm]
	if !ok {
		return nil, fmt.Errorf("invalid BanMatchingType: %d", bm)
	}
	return s, nil
}

func (bm *BanMatchingType) UnmarshalYAML(value *yaml.Node) error {
	v, ok := _BanMatchingTypeNameToValue[value.Value]
	if !ok {
		return fmt.Errorf("invalid BanMatchingType %q", value.Value)
	}
	*bm = v
	return nil
}

func (bm BanMatchingType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(bm).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _BanMatchingTypeValueToName[bm]
	if !ok {
		return nil, fmt.Errorf("invalid BanMatchingType: %d", bm)
	}
	return json.Marshal(s)
}

func (bm *BanMatchingType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("BanMatchingType should be a string, got %s", data)
	}
	v, ok := _BanMatchingTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid BanMatchingType %q", s)
	}
	*bm = v
	return nil
}

func (bm BanMatchingType) Val() int {
	return int(bm)
}

func (bm BanMatchingType) String() string {
	return _BanMatchingTypeValueToName[bm]
}
