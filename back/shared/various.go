package shared

import (
    "strconv"
    "time"
    "github.com/golang/protobuf/ptypes"
    "github.com/golang/protobuf/ptypes/timestamp"
)

func Uint32Min(a, b uint32) uint32 {
    if a < b {
        return a
    }
    return b
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
