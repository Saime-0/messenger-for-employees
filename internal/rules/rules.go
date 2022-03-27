package rules

const (
	NameMaxLength      = 32
	NameMinLength      = 1
	DomainMaxLength    = 32
	DomainMinLength    = 4
	NoteMaxLength      = 64
	RefreshTokenLength = 28
	MaxLimit           = 20
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
)

type AdvancedError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (p *AdvancedError) Error() string {
	return p.Message
}
