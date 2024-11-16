package internal

import "time"

type Time struct {
	time.Time
}

func (ct *Time) UnmarshalCSV(csvField string) error {
	const layout = "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, csvField)
	if err != nil {
		return err
	}
	ct.Time = parsedTime

	return nil
}
