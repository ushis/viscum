package util

import (
  "code.google.com/p/goconf/conf"
)

// Default configuration
var defaults = map[string]map[string]string{
  "database": map[string]string{
    "auth":   "dbname=viscum user=viscum password=secret host=localhost port5432",
    "driver": "postgres",
  },
  "mailer": map[string]string{
    "command": "/usr/bin/mail",
  },
  "rpc": map[string]string{
    "socket": "/run/viscum/viscum.sock",
  },
}

type Config struct {
  file *conf.ConfigFile // Config file.
}

// Reads the config file.
func ReadConfig(name string) (*Config, error) {
  file, err := conf.ReadConfigFile(name)

  if err != nil {
    return nil, err
  }
  return &Config{file: file}, nil
}

// Returns a config value specified by section and key.
func (self *Config) Get(sec string, key string) string {
  if val, err := self.file.GetString(sec, key); err == nil {
    return val
  }
  if _, ok := defaults[sec]; !ok {
    Fatal("Section not found:", sec)
  }
  if val, ok := defaults[sec][key]; ok {
    return val
  }
  Fatal("Value not found:", sec, key)
  return "" // Never...
}