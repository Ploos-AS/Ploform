package main

import (
 "errors"
 "fmt"
 "net/url"
 "os"
 "strings"
)

// databaseDSN returns an explicit URL if supplied, otherwise constructs a
// local Docker-network connection string from a mounted password file.
func databaseDSN(explicitEnv, passwordFileEnv, username string) (string, error) {
 if dsn := os.Getenv(explicitEnv); dsn != "" { return dsn, nil }
 path := os.Getenv(passwordFileEnv)
 if path == "" { return "", nil }
 raw, err := os.ReadFile(path)
 if err != nil { return "", fmt.Errorf("read database secret: %w", err) }
 password := strings.TrimSuffix(string(raw), "\n")
 if len(password) < 24 || strings.ContainsAny(password, "\r\n") {
  return "", errors.New("invalid database password file")
 }
 host := os.Getenv("PLOFORM_DATABASE_HOST")
 if host == "" { host = "postgres" }
 dbname := os.Getenv("PLOFORM_DATABASE_NAME")
 if dbname == "" { dbname = "ploform" }
 if strings.ContainsAny(host, "/?#@") || strings.ContainsAny(dbname, "/?#@") {
  return "", errors.New("invalid database host or name")
 }
 u := &url.URL{Scheme: "postgres", User: url.UserPassword(username, password), Host: host + ":5432", Path: "/" + dbname}
 q := u.Query()
 q.Set("sslmode", "disable") // only on an isolated local development Docker network
 u.RawQuery = q.Encode()
 return u.String(), nil
}
