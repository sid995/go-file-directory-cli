package times

import (
	"fmt"
	"time"
)

func FormatRelativeTime(modTime time.Time) string {
	today := time.Now().Truncate(24 * time.Hour)
	modTime = modTime.Truncate(24 * time.Hour)

	duration := today.Sub(modTime)

	if duration.Hours() < 24 {
		return fmt.Sprintf("%s ago", duration.Round(time.Minute).String())
	}

	return modTime.Format(time.RFC3339)
}
