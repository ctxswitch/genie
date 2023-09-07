package resources

import (
	"ctx.sh/genie/pkg/resources/integer_range"
	"ctx.sh/genie/pkg/resources/list"
	"ctx.sh/genie/pkg/resources/random_string"
	"ctx.sh/genie/pkg/resources/timestamp"
	"ctx.sh/genie/pkg/resources/uuid"
)

type Config struct {
	IntegerRanges map[string]integer_range.Config `yaml:"integer_ranges"`
	Lists         map[string]list.Config          `yaml:"lists"`
	RandomStrings map[string]random_string.Config `yaml:"random_strings"`
	Timestamps    map[string]timestamp.Config     `yaml:"timestamps"`
	Uuids         map[string]uuid.Config          `yaml:"uuids"`
}
