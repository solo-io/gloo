package constants

const (
	// DiscardInterruptedTranslationEnv turns on handling of translation and envoy validation that was
	// interrupted by a cancelled context or a killed envoy validation process (default off):
	//   - An interrupted envoy validation is not reported or cached as a config error, and the
	//     translation that ran it is discarded instead of published.
	//   - A gateway proxy reconcile or xDS sync whose context ended does not write Proxies or xDS.
	//   - A gloo setup restart waits until the running setup's first translator Sync returns, so each
	//     setup publishes xDS at least once.
	DiscardInterruptedTranslationEnv = "GLOO_DISCARD_INTERRUPTED_TRANSLATION"
)
