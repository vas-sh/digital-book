package config

var Config = struct {
	DSN string
}{
	DSN: "host=localhost user=user password=1111 dbname=test port=5432 sslmode=disable TimeZone=Europe/Kiev",
}
