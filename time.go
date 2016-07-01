package pipedrive

import "time"

type Time struct {
	time.Time
}

func NewTime(t time.Time) Time {
	return Time{t}
}

func (t *Time) UnmarshalJSON(buf []byte) error {
	//fmt.Println(string(buf))
	tt, err := time.Parse("2006-01-02 15:04:05", string(buf[1:len(buf)-1]))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

func (t *Time) String() string {
	if t == nil {
		return ""
	}
	return t.Time.Local().Format("2006-01-02 15:04:05")
}

type Date struct {
	time.Time
}

func NewDate(t time.Time) Date {
	return Date{t}
}

func (t *Date) UnmarshalJSON(buf []byte) error {
	//fmt.Println(string(buf))
	tt, err := time.Parse("2006-01-02", string(buf[1:len(buf)-1]))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

func (t *Date) String() string {
	if t == nil {
		return ""
	}
	return t.Time.Local().Format("2006-01-02")
}
