package enums

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type BannRuleType int8

const (
	NETWORK BannRuleType = iota + 1
	MAC
	IP
	HOST
	PORT
)

var (
	_BannRuleTypeNameToValue = map[string]BannRuleType{
		"NETWORK": NETWORK,
		"network": NETWORK,
		"NET":     NETWORK,
		"net":     NETWORK,
		"MAC":     MAC,
		"mac":     MAC,
		"IP":      IP,
		"ip":      IP,
		"HOST":    HOST,
		"host":    HOST,
		"PORT":    PORT,
		"port":    PORT,
	}

	_BannRuleTypeValueToName = map[BannRuleType]string{
		NETWORK: "net",
		MAC:     "mac",
		IP:      "ip",
		HOST:    "host",
		PORT:    "port",
	}
)

func (br BannRuleType) MarshalYAML() (interface{}, error) {
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
