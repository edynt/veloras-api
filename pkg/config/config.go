package config

type Config struct {
	PostgreSQL PostgreSQLSetting `mapstructure:"postgresql"`
	Logger     LoggerSetting     `mapstructure:"logger"`
	Server     ServerSetting     `mapstructure:"server"`
	JWT        JWTSetting        `mapstructure:"jwt"`
	SMTP       SMTPSetting       `mapstructure:"mail"`
	Cron       CronSetting       `mapstructure:"cron"`
}

type LoggerSetting struct {
	Log_level     string `mapstructure:"log_level"`
	File_log_name string `mapstructure:"file_log_name"`
	Max_size      int    `mapstructure:"max_size"`
	Max_backups   int    `mapstructure:"max_backups"`
	Max_age       int    `mapstructure:"max_age"`
	Compress      bool   `mapstructure:"compress"`
}

type ServerSetting struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type PostgreSQLSetting struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	ConnMaxLifeTime int    `mapstructure:"connMaxLifeTime"`
	SslMode         string `mapstructure:"sslmode"`
}

type JWTSetting struct {
	TokenHourLifespan  uint   `mapstructure:"token_hour_lifespan"`
	ApiSecret          string `mapstructure:"api_secret"`
	AccessTokenExpire  int    `mapstructure:"access_token_expire"`
	RefreshTokenExpire int    `mapstructure:"refresh_token_expire"`
}

type SMTPSetting struct {
	Host     string `mapstructure:"smtp_host"`
	Port     string `mapstructure:"smtp_port"`
	User     string `mapstructure:"smtp_user"`
	Password string `mapstructure:"smtp_password"`
}

type CronSetting struct {
	Enabled               bool `mapstructure:"enabled"`
	CleanupSessions       bool `mapstructure:"cleanup_sessions"`
	CleanupTokens         bool `mapstructure:"cleanup_tokens"`
	CleanupVerifications  bool `mapstructure:"cleanup_verifications"`
	CleanupPasswordResets bool `mapstructure:"cleanup_password_resets"`
}
