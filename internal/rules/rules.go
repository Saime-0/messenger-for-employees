package rules

const (
	RefreshTokenLength = 28
	DefaultLimitValue  = 50 // сделать методы в node Для каждого типа запроса и разделить валидацию и заполнение стандартными значениями
	MaxMsgCount        = DefaultLimitValue
	MinPasswordLength  = 6
	MaxPasswordLength  = 32
	//RefreshTokenLifetime          = int64(60 * 60 * 24 * 60) // 60 days
	MaxRefreshSession             = 5
	LifetimeOfMarkedClient        = int64(60)      // s.
	LifetimeOfRegistrationSession = int64(60 * 60) // 1 hour
	DurationOfScheduleInterval    = int64(60)      // 1 hour

	AllowedConnectionShutdownDuration = 120

	MaxFirstNameLen = 32
	MaxLastNameLen  = 32
	MaxFullNameLen  = MaxFirstNameLen + MaxLastNameLen + 1 // 1 = space

	MaxMessageBodyLen = 1024
)

type AdvancedError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (p *AdvancedError) Error() string {
	return p.Message
}
