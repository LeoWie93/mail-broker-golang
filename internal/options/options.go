package options

import "fmt"

type Options struct {
	Template string
}

var optionsMap = map[string]*Options{
	"v1": &Options{Template: "v1.html"},
	"v2": &Options{Template: "v2.html"},
}

func ExchangeAction(action string) (options *Options, err error) {
	if options, ok := optionsMap[action]; ok {
		return options, nil
	}

	return nil, fmt.Errorf("Given action is not valid: %s", action)
}
