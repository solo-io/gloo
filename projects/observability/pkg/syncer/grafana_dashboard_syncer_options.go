package syncer

type grafanaDashboardsSyncerOpts struct {
	defaultDashboardFolderId   uint
	dashboardPrefix            string
	extraMetricQueryParameters string
	extraDashboardTags         []string
}

type Option func(*grafanaDashboardsSyncerOpts)

func WithDefaultDashboardFolderId(folderID uint) Option {
	return func(opts *grafanaDashboardsSyncerOpts) {
		opts.defaultDashboardFolderId = folderID
	}
}

func WithDashboardPrefix(prefix string) Option {
	return func(opts *grafanaDashboardsSyncerOpts) {
		opts.dashboardPrefix = prefix
	}
}

func WithExtraMetricQueryParameters(extraMetricQueryParameters string) Option {
	return func(opts *grafanaDashboardsSyncerOpts) {
		opts.extraMetricQueryParameters = extraMetricQueryParameters
	}
}

func WithExtraDashboardTags(extraDashboardTags []string) Option {
	return func(opts *grafanaDashboardsSyncerOpts) {
		opts.extraDashboardTags = extraDashboardTags
	}
}

func processOptions(options ...Option) grafanaDashboardsSyncerOpts {
	// Set defaults here
	opts := &grafanaDashboardsSyncerOpts{
		defaultDashboardFolderId: generalFolderId,
		// Set it to Avoid NPE when referring to this field
		extraDashboardTags: []string{},
	}
	for _, o := range options {
		o(opts)
	}
	return *opts
}
