// Package filter provides composable log line filters for logslice.
//
// # RegionFilter
//
// RegionFilter matches log lines where a specified field contains one of
// the provided cloud region identifiers. Matching is case-insensitive.
//
// Patterns may include a trailing wildcard "*" to match all regions
// sharing a given prefix:
//
//	// Match exactly us-east-1 or eu-west-2
//	f, _ := filter.NewRegionFilter("region", []string{"us-east-1", "eu-west-2"})
//
//	// Match any US region (us-east-1, us-west-2, us-gov-west-1, …)
//	f, _ := filter.NewRegionFilter("region", []string{"us-*"})
//
// CLI usage (--region flag):
//
//	logslice --region region:us-east-1,eu-west-2 app.log
//	logslice --region dc:us-* app.log
package filter
