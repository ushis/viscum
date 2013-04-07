package util

import (
  "code.google.com/p/goconf/conf"
  "os"
  "time"
)

// Default configuration
var defaults = map[string]map[string]string{
  "database": map[string]string{
    "auth":   "dbname=viscum user=viscum password=secret host=localhost port=5432",
    "driver": "postgres",
  },
  "feed": map[string]string{
    "poll": "15m",
  },
  "mail": map[string]string{
    "from":          "viscum@localhost",
    "mailer":        "pipe",
    "pipe":          "/usr/bin/mail",
    "smtp_host":     "127.0.0.1",
    "smtp_port":     "25",
    "smtp_username": "viscum",
    "smtp_password": "secret",
  },
  "rpc": map[string]string{
    "socket": "/run/viscum/viscum.sock",
  },
}

type Config struct {
  *conf.ConfigFile
}

// Reads the config file.
func ReadConfig(name string) (*Config, error) {
  file, err := os.Open(name)

  if err != nil {
    return nil, err
  }
  defer file.Close()

  c := &Config{conf.NewConfigFile()}

  for sec, opt := range defaults {
    for k, v := range opt {
      c.AddOption(sec, k, v)
    }
  }

  if err = c.Read(file); err != nil {
    return nil, err
  }
  return c, nil
}

// Returns a config value specified by section and key.
func (self *Config) Get(sec string, key string) string {
  val, err := self.GetString(sec, key)

  if err != nil {
    Fatal(err)
  }
  return val
}

// Returns a duration.
func (self *Config) GetDuration(sec string, key string) time.Duration {
  val, err := time.ParseDuration(self.Get(sec, key))

  if err != nil {
    Fatal(err)
  }
  return val
}
