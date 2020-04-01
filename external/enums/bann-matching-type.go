package enums

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type BannMatchingType int8

const (
	Strict BannMatchingType = iota + 1
	ByPattern
)

var (
	_BannMatchingTypeNameToValue = map[string]BannMatchingType{
		"Strict":     Strict,
		"ByPattern":  ByPattern,
		"strict":     Strict,
		"by_pattern": ByPattern,
	}

	_BannMatchingTypeValueToName = map[BannMatchingType]string{
		ByPattern: "ByPattern",
		Strict:    "Strict",
	}
)

func (br BannMatchingType) MarshalYAML() (interface{}, error) {
	s, ok := _BannRuleTypeValueToName[br]
	if !ok {
		return nil, fmt.Errorf("invalid BannRuleType: %d", br)
	}
	return s, nil
}

func (br *BannRuleType) UnmarshalYAML(value *yaml.Node) error {
	v, ok := _BannRuleTypeNameToValue[value.Value]
	if !ok {
		return fmt.Errorf("invalid BannRuleType %q", value.Value)
	}
	*br = v
	return nil
}

func (br BannRuleType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(br).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _BannRuleTypeValueToName[br]
	if !ok {
		return nil, fmt.Errorf("invalid BannRuleType: %d", br)
	}
	return json.Marshal(s)
}

func (br *BannRuleType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("BannRuleType should be a string, got %s", data)
	}
	v, ok := _BannRuleTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid BannRuleType %q", s)
	}
	*br = v
	return nil
}

func (br BannRuleType) Val() int {
	return int(br)
}

func (br BannRuleType) String() string {
	return _BannRuleTypeValueToName[br]
}
