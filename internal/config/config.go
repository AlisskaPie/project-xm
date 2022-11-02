package config

type Config struct {
	DB          DB
	HTTP        HTTP
	Auth        Auth
	EventSender bool
}

type DB struct {
	DSN string
}

type HTTP struct {
	ListenHostPort string
}

type Auth struct {
	JWTKey string
}
