package times

import "time"

var shanghai, _ = time.LoadLocation("Asia/Shanghai")

func In(value time.Time) time.Time {
	return value.In(shanghai)
}

func Parse(layout, value string) (time.Time, error) {
	if time.RFC3339 == layout {
		result, err := time.Parse(layout, value)
		if err != nil {
			return result, err
		}
		return In(result), nil
	}
	result, err := time.ParseInLocation(layout, value, shanghai)
	if err != nil {
		return result, err
	}
	return result, nil
}
