package tweets
import (
	"time"
)

const ctLayout = "Mon Jan 02 15:04:05 -0700 2006"

type CustomTime struct {
    time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
    if b[0] == '"' && b[len(b)-1] == '"' {
        b = b[1 : len(b)-1]
    }
    ct.Time, err = time.Parse(ctLayout, string(b))
    return
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
    return []byte(ct.Time.Format(ctLayout)), nil
}

var nilTime = (time.Time{}).UnixNano()
func (ct *CustomTime) IsSet() bool {
    return ct.UnixNano() != nilTime
}
