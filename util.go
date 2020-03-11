package main

import (
	"fmt"
	"time"
)


// MarshalJSON 发出适合在json中使用的时间戳
func (t *Timestamp) MarshalJSON() ([]byte, error) {

	fmt.Print()
	ts := time.Time(*t).Format(time.RFC3339)
	stamp := fmt.Sprint(ts)
	return []byte(stamp), nil
}

// UnmarshalJSON 从json输入解析时间戳
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]

	fmt.Println(s)

	ts, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		ts, err = time.Parse("2006-01-02T15:04:05Z07:00", s)
		if err != nil {
			return err
		}
	}
	*t = Timestamp(ts)
	return nil
}
