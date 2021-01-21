package lightmq

// Options ...
type Options struct {
	// Hostname of where broker should listen
	//
	// Default: "0.0.0.0"
	Hostname string

	// Port of where broker should listen
	//
	// Default: "1883"
	Port uint32
}

// Parse parses options and set defaults
func (opts Options) Parse() Options {
	if opts.Hostname == "" {
		opts.Hostname = "0.0.0.0"
	}
	if opts.Port == 0 {
		opts.Port = 1883
	}

	return opts
}
