package template

type templateGeneratorOpts struct {
	dashboardPrefix            string
	extraMetricQueryParameters string
	extraDashboardTags         []string
}

type Option func(*templateGeneratorOpts)

func WithDashboardPrefix(prefix string) Option {
	return func(opts *templateGeneratorOpts) {
		opts.dashboardPrefix = prefix
	}
}

func WithExtraMetricQueryParameters(extraMetricQueryParameters string) Option {
	return func(opts *templateGeneratorOpts) {
		opts.extraMetricQueryParameters = extraMetricQueryParameters
	}
}

func WithExtraDashboardTags(extraDashboardTags []string) Option {
	return func(opts *templateGeneratorOpts) {
		opts.extraDashboardTags = extraDashboardTags
	}
}

func processOptions(options ...Option) templateGeneratorOpts {
	opts := &templateGeneratorOpts{
		// Set it to Avoid NPE when referring to this field
		extraDashboardTags: []string{},
	}
	for _, o := range options {
		o(opts)
	}
	return *opts
}
