package service

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// toTimestamp 时间 → proto Timestamp（零值也转）。
func toTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

// toTimestampOrNil 可空时间：零值表示"没有"，转成 nil。
func toTimestampOrNil(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

// tsTime proto Timestamp → 时间，nil 返回零值。
func tsTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
