package units

import "time"


// Timespan is a helper type for formatting of laptimes.
type Timespan time.Duration

// Format formats the Timespan as a prettified laptime. 
func (t Timespan) Format(format string) string {
	return time.Unix(0, 0).UTC().Add(time.Duration(t)).Format(format)
}
