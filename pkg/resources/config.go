package resources

import (
	"ctx.sh/genie/pkg/resources/float_range"
	"ctx.sh/genie/pkg/resources/integer_range"
	"ctx.sh/genie/pkg/resources/ipaddr"
	"ctx.sh/genie/pkg/resources/list"
	"ctx.sh/genie/pkg/resources/random_string"
	"ctx.sh/genie/pkg/resources/timestamp"
	"ctx.sh/genie/pkg/resources/uuid"
)

// Config is the configuration for a collection of resources.  It is the
// top-level configuration for the resources block.
type Config struct {
	IntegerRanges map[string]integer_range.Config `yaml:"integer_ranges"`
	FloatRanges   map[string]float_range.Config   `yaml:"float_ranges"`
	Lists         map[string]list.Config          `yaml:"lists"`
	RandomStrings map[string]random_string.Config `yaml:"random_strings"`
	Timestamps    map[string]timestamp.Config     `yaml:"timestamps"`
	UUIDs         map[string]uuid.Config          `yaml:"uuids"`
	IPAddrs       map[string]ipaddr.Config        `yaml:"ipaddrs"`
}
