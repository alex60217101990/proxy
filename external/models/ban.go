package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/alex60217101990/proxy.git/external/enums"
)

type FirewallConfigs struct {
	// ... Timeout for studio
}

type BlockObject struct {
	RuleType     enums.BanRuleType     `yaml:"rule_type" json:"rule_type"`
	MatchingType enums.BanMatchingType `yaml:"matching_type" json:"matching_type"`
	RuleBody     string                `yaml:"rule_body" json:"rule_body"`
}

type BanEvent struct {
	Name string
	IP   string
}

type BanObject struct {
	Name        string        `yaml:"rule_type" json:"rule_type"`
	StrikeLimit uint16        `yaml:"strike_limit" json:"strike_limit"`
	ExpireBase  time.Duration `yaml:"expire_base" json:"expire_base"`
	Sentence    time.Duration `yaml:"sentence" json:"sentence"`
}

func (b *BanObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		Name        string `yaml:"rule_type"`
		StrikeLimit uint16 `yaml:"strike_limit"`
		ExpireBase  string `yaml:"expire_base"`
		Sentence    string `yaml:"sentence"`
	}

	var tmp alias
	if err := unmarshal(&tmp); err != nil {
		return err
	}
	b.Name = tmp.Name

	t, err := time.ParseDuration(tmp.ExpireBase)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' to time.Duration: %v", tmp.ExpireBase, err)
	}
	b.ExpireBase = t

	t, err = time.ParseDuration(tmp.Sentence)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' to time.Duration: %v", tmp.Sentence, err)
	}
	b.Sentence = t

	return nil
}

func (b *BanObject) UnmarshalJSON(bts []byte) error {
	type alias struct {
		Name        string `json:"rule_type"`
		StrikeLimit uint16 `json:"strike_limit"`
		ExpireBase  string `json:"expire_base"`
		Sentence    string `json:"sentence"`
	}

	var tmp alias
	if err := json.Unmarshal(bts, &tmp); err != nil {
		return err
	}
	b.Name = tmp.Name

	t, err := time.ParseDuration(tmp.ExpireBase)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' to time.Duration: %v", tmp.ExpireBase, err)
	}
	b.ExpireBase = t

	t, err = time.ParseDuration(tmp.Sentence)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' to time.Duration: %v", tmp.Sentence, err)
	}
	b.Sentence = t

	return nil
}
