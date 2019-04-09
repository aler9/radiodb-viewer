package shared

import (
    "fmt"
    "strconv"
    "strings"
    "unicode"
    "time"
    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
    "github.com/golang/protobuf/ptypes"
    "github.com/golang/protobuf/ptypes/timestamp"
)

func Uint32Min(a, b uint32) uint32 {
    if a < b {
        return a
    }
    return b
}

func isMn(r rune) bool {
    return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func GetTextKeywords(in string, minLength int) (map[string]struct{}) {
    // replace accents
    t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
    in,_,_ = transform.String(t, in)

    // to lower
    in = strings.ToLower(in)

    // split by space
    ret := make(map[string]struct{})
    for _,word := range strings.Split(in, " ") {
        if len(word) < minLength { // remove too short
            continue
        }
        ret[word] = struct{}{}
    }
    return ret
}

func FormatFirstSeen(ts *timestamp.Timestamp, format string) string {
    date := PbtimeToTime(ts)

    ret := ""
    cmp := time.Date(2018, time.September, 6, 24, 59, 59, 0, time.UTC)
    if date.Before(cmp) == true {
        ret += "before "
    }
    ret += date.Format(format)
    return ret
}

func FormatDuration(duration float64) string {
    d := time.Duration(duration * 1000.0) * time.Millisecond
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    d -= m * time.Minute
    s := d / time.Second

    var ret []string
    if h > 0 {
        ret = append(ret, fmt.Sprintf("%dh", h))
    }
    if m > 0 {
        ret = append(ret, fmt.Sprintf("%dm", m))
    }
    ret = append(ret, fmt.Sprintf("%ds", s))
    return strings.Join(ret, " ")
}

func Atoui32(s string) uint32 {
    ret,_ := strconv.ParseUint(s, 10, 32)
    return uint32(ret)
}

func TimeToPbtime(in time.Time) *timestamp.Timestamp {
    ret,err := ptypes.TimestampProto(in)
    if err != nil {
        panic(err)
    }
    return ret
}

func PbtimeToTime(in *timestamp.Timestamp) time.Time {
    ret,err := ptypes.Timestamp(in)
    if err != nil {
        panic(err)
    }
    return ret
}
