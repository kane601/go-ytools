package tools

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ToHumanizeDurtion(dur time.Duration) (hum string) {
	h := int(dur / time.Hour)
	m := int((dur - time.Duration(h)*time.Hour) / time.Minute)
	s := int((dur - time.Duration(h)*time.Hour - time.Duration(m)*time.Minute) / time.Second)
	if h > 0 {
		hum = fmt.Sprintf("%d:%02d:%02d", h, m, s)
	} else {
		hum = fmt.Sprintf("%d:%02d", m, s)
	}
	return hum
}

func ParseHumanizeDuration(hum string) (time.Duration, error) {
	parts := strings.Split(hum, ":")
	if len(parts) == 2 {
		// Format: MM:SS
		m, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		s, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		return time.Duration(m)*time.Minute + time.Duration(s)*time.Second, nil
	} else if len(parts) == 3 {
		// Format: HH:MM:SS
		h, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		s, err := strconv.Atoi(parts[2])
		if err != nil {
			return 0, err
		}
		return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second, nil
	} else {
		return 0, fmt.Errorf("invalid format")
	}
}
