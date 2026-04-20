// Package filter provides composable log-line filters for logslice.
//
// # TagFilter
//
// TagFilter matches log lines where a named field contains at least one of a
// set of user-supplied tags.  The field value is split by a configurable
// separator (default ",") and each segment is compared case-insensitively.
//
// Example — match lines whose "tags" field includes "prod" or "canary":
//
//	f, err := filter.NewTagFilter("tags", ",", []string{"prod", "canary"})
//	if err != nil { ... }
//	matched := f.Match(line)
//
// CLI flag format (via ParseTagFlag):
//
//	--tag 'tags:prod,canary'
//
// A custom separator can be appended as a third colon-delimited segment:
//
//	--tag 'tags:prod|canary:|'
package filter
