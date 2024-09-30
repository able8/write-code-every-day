package main

import (
	"log"
	"log/slog"
	"os"
)

type User struct {
	ID       int
	Username string
	Password string
}

// Currently, sensitive data can be leaked out to the log.
//  If we implement the slog.LogValuer interface,
// we can prevent our sensitive data from leaking out to the logs.

func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("ID", u.ID),
		slog.String("username", u.Username),
		slog.String("password", "TOKEN_REDACTED"),
	)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	log.Println("program started")

	u := User{ID: 10, Username: "admin", Password: "abc123"}
	// the log package does NOT detect the `LogValuer` interface
	log.Println(u)

	// You must use the slog package to detect the `Log.Valuer`
	slog.Info("user", "user", u)

}
