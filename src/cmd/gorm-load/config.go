package main

var cfg *Config

type Config struct {
	Database ConfigDatabase `json:"database" yaml:"database" toml:"database" xml:"database" ini:"database"`
}

type ConfigDatabase struct {
	ActiveProfile string                  `json:"active_profile" yaml:"active_profile" toml:"active_profile" xml:"active_profile" ini:"active_profile"`
	Profiles      []ConfigDatabaseProfile `json:"profiles" yaml:"profiles" toml:"profiles" xml:"profiles" ini:"profiles"`
}

type ConfigDatabaseProfile struct {
	Database     string `json:"database" yaml:"database" toml:"database" xml:"database" ini:"database"`
	Hostname     string `json:"hostname" yaml:"hostname" toml:"hostname" xml:"hostname" ini:"hostname"`
	ID           string `json:"id" yaml:"id" toml:"id" xml:"id" ini:"id"`
	Insecure     bool   `json:"insecure" yaml:"insecure" toml:"insecure" xml:"insecure" ini:"insecure"`
	Password     string `json:"password" yaml:"password" toml:"password" xml:"password" ini:"password"`
	Port         int    `json:"port" yaml:"port" toml:"port" xml:"port" ini:"port"`
	RootPassword string `json:"root_password" yaml:"root_password" toml:"root_password" xml:"root_password" ini:"root_password"`
	User         string `json:"user" yaml:"user" toml:"user" xml:"user" ini:"user"`
	UserPassword string `json:"user_password" yaml:"user_password" toml:"user_password" xml:"user_password" ini:"user_password"`
}
