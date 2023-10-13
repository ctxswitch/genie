package resources

import (
	"stvz.io/genie/pkg/resources/integer_range"
	"stvz.io/genie/pkg/resources/ipaddr"
	"stvz.io/genie/pkg/resources/list"
	"stvz.io/genie/pkg/resources/random_string"
	"stvz.io/genie/pkg/resources/timestamp"
	"stvz.io/genie/pkg/resources/uuid"
)

type Config struct {
	IntegerRanges map[string]integer_range.Config `yaml:"integer_ranges"`
	Lists         map[string]list.Config          `yaml:"lists"`
	RandomStrings map[string]random_string.Config `yaml:"random_strings"`
	Timestamps    map[string]timestamp.Config     `yaml:"timestamps"`
	Uuids         map[string]uuid.Config          `yaml:"uuids"`
	IPAddrs       map[string]ipaddr.Config        `yaml:"ipaddrs"`
}
