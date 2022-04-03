package res

// regular
const (
	Year       int64 = 31536000
	AuthHeader       = "Authorization"
)

type FetchType string

const (
	Positive FetchType = "POSITIVE"
	Neutral  FetchType = "NEUTRAL"
	Negative FetchType = "NEGATIVE"
)

type UnitType string

const (
	User UnitType = "USER"
	Chat UnitType = "CHAT"
)

type LocalKeys int

const (
	_ LocalKeys = iota

	// ctx keys
	CtxAuthData
	CtxUserAgent
	CtxNode
	CtxAdminData

	// cache keys
	CacheNextRunRegularScheduleAt
	CacheCurrentReconnectionAttemptToLogDB
	CacheScheduleInvites

	// indicators
	IndicatorLogger
	// states
	OK
	FailedDBConnection
	RepairingConnection
)
