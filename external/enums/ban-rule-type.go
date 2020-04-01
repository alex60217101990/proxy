package enums

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type BanRuleType int8

const (
	NETWORK BanRuleType = iota + 1
	MAC
	IP
	HOST
	PORT
)

var (
	_BanRuleTypeNameToValue = map[string]BanRuleType{
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

	_BanRuleTypeValueToName = map[BanRuleType]string{
		NETWORK: "net",
		MAC:     "mac",
		IP:      "ip",
		HOST:    "host",
		PORT:    "port",
	}
)

func (br BanRuleType) MarshalYAML() (interface{}, error) {
	s, ok := _BanRuleTypeValueToName[br]
	if !ok {
		return nil, fmt.Errorf("invalid BanRuleType: %d", br)
	}
	return s, nil
}

func (br *BanRuleType) UnmarshalYAML(value *yaml.Node) error {
	v, ok := _BanRuleTypeNameToValue[value.Value]
	if !ok {
		return fmt.Errorf("invalid BanRuleType %q", value.Value)
	}
	*br = v
	return nil
}

func (br BanRuleType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(br).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _BanRuleTypeValueToName[br]
	if !ok {
		return nil, fmt.Errorf("invalid BanRuleType: %d", br)
	}
	return json.Marshal(s)
}

func (br *BanRuleType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("BanRuleType should be a string, got %s", data)
	}
	v, ok := _BanRuleTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid BanRuleType %q", s)
	}
	*br = v
	return nil
}

func (br BanRuleType) Val() int {
	return int(br)
}

func (br BanRuleType) String() string {
	return _BanRuleTypeValueToName[br]
}
