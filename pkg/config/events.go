package config

type EventOptions struct {
	Tags     map[string]string `yaml:"tags"`
	Source   string            `yaml:"source"`
	Template string            `yaml:"template"`
	Payload  string            `yaml:"payload"`
}

// func (e *EventOptions) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	type RawEventOptions EventOptions

// 	var str string
// 	if err := unmarshal(&str); err == nil {
// 		e.Raw = str
// 		return nil
// 	}

// 	var out RawEventOptions
// 	if err := unmarshal(&out); err != nil {
// 		return err
// 	}

// 	*e = EventOptions(out)
// 	return nil
// }

type Event struct {
	Count string            `yaml:"count"`
	Vars  map[string]string `yaml:"vars"`
	Event EventOptions      `yaml:"event"`
}

// UnmarshalYAML sets the Event defaults and parses an event block. The
// event block can either be a string or a transport event.
func (e *Event) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type EventDefaulted Event
	var defaults = EventDefaulted{
		Count: "1",
		Vars:  make(map[string]string),
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	*e = Event(out)
	return nil
}
