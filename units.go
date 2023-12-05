package wrc

import "time"

// # Conversion values
const (
	
	MpsToMph  float32 = 2.237
	MpsToKmph float32 = 3.6
)

// Timespan is a helper type used to format laptimes.
//
// 	Timespan(time.Duration(10 * time.Second)).Format("04:05.000")
//   	"00:10.000"
type Timespan time.Duration

// Format formats the Timespan as a prettified laptime.
func (t Timespan) Format(format string) string {
	return time.Unix(0, 0).UTC().Add(time.Duration(t)).Format(format)
}
