package config

import (
	"github.com/BurntSushi/toml"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/clog"
	"os"
)

type Config2 struct {
	FromEnv
	FromCfgFile
}

func NewConfig2(pathToCfgFile string) (*Config2, error) {
	fromFile := new(FromCfgFile)
	_, err := toml.DecodeFile(pathToCfgFile, fromFile)
	if err != nil {
		return nil, cerrors.Wrap(err, "не удалось декодировать файл")
	}

	if !fromFile.validate() {
		return nil, cerrors.New("в файле конфигурации заполнены не все поля")
	}
	if !clog.Exists(clog.LogLevel(*fromFile.Logging.LoggingLevel)) {
		return nil, cerrors.New("указан несуществующий уровень лога")
	}

	fromEnv := &FromEnv{
		PostgresConnection: os.Getenv("POSTGRES_CONNECTION"),
		GlobalPasswordSalt: os.Getenv("GLOBAL_PASSWORD_SALT"),
		MongoDBUri:         os.Getenv("MONGODB_URI"),
		SecretSigningKey:   os.Getenv("SECRET_SIGNING_KEY"),
		SmtpHost:           os.Getenv("SMTP_HOST"),
		SmtpEmailLogin:     os.Getenv("SMTP_EMAIL_LOGIN"),
		SmtpEmailPasswd:    os.Getenv("SMTP_EMAIL_PASSWD"),
	}

	if !fromEnv.validate() {
		return nil, cerrors.New("не установлены некоторые переменные окружения")
	}
	return &Config2{
		FromEnv:     *fromEnv,
		FromCfgFile: *fromFile,
	}, nil
}

type FromEnv struct {
	PostgresConnection string // `toml:"postgres_connection"`
	GlobalPasswordSalt string // `toml:"global_password_salt"`
	MongoDBUri         string // `toml:"mongodb_uri"`
	SecretSigningKey   string // `toml:"secret_signing_key"`
	SmtpHost           string // `toml:"smtp_host"`
	SmtpEmailLogin     string // `toml:"smtp_email_login"`
	SmtpEmailPasswd    string // `toml:"smtp_email_passwd"`
}

type FromCfgFile struct {
	ApplicationPort                   *string  `toml:"application_port" json:"application_port,omitempty"`
	QueryComplexityLimit              *int     `toml:"query_complexity_limit" json:"query_complexity_limit,omitempty"`
	DurationOfScheduleInterval        *int64   `toml:"duration_of_schedule_interval" json:"duration_of_schedule_interval,omitempty"`
	RefreshTokenLifetime              *int64   `toml:"refresh_token_lifetime" json:"refresh_token_lifetime,omitempty"`
	AccessTokenLifetime               *int64   `toml:"access_token_lifetime" json:"access_token_lifetime,omitempty"`
	MaximumNumberOfMessagesPerRequest *int     `toml:"maximum_number_of_messages_per_request" json:"maximum_number_of_messages_per_request,omitempty"`
	MaxCountRooms                     *int     `toml:"max_count_rooms" json:"max_count_rooms,omitempty"`
	MaxUserChats                      *int     `toml:"max_user_chats" json:"max_user_chats,omitempty"`
	MaxCountOwnedChats                *int     `toml:"max_count_owned_chats" json:"max_count_owned_chats,omitempty"`
	MaxMembersOnChat                  *int     `toml:"max_members_on_chat" json:"max_members_on_chat,omitempty"`
	MaxRolesInChat                    *int     `toml:"max_roles_in_chat" json:"max_roles_in_chat,omitempty"`
	SMTPing                           *SMTPing `toml:"smtp" json:"smt_ping,omitempty"`
	Logging                           *Logging `toml:"log" json:"logging,omitempty"`
}

type SMTPing struct {
	SMTPAuthor *string `toml:"smtp_author"`
	SMTPPort   *int    `toml:"smtp_port"`
}

type Logging struct {
	LoggingOutput *uint8  `toml:"logging_output"`
	LoggingLevel  *int8   `toml:"logging_level"`
	LoggingDBName *string `toml:"logging_db_name"`
}
