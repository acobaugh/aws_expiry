package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-ini/ini"
)

const key = "x_security_token_expires"

func main() {
	fTimeLeft := flag.Bool("time-left", false, "Show time remaining")
	fErrorOnExpired := flag.Bool("error-on-expired", false, "Exit 1 if creds are expired")
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error: %s\n", homeDir)
		os.Exit(1)
	}

	profile := os.Getenv("AWS_PROFILE")
	sharedCredsFile := fmt.Sprintf("%s/%s", homeDir, ".aws/credentials")

	i, err := ini.Load(sharedCredsFile)
	handleError(err)

	s, err := i.GetSection(profile)
	handleError(err)

	k, err := s.GetKey(key)
	handleError(err)

	// Mon Jan 2 15:04:05 MST 2006
	// 2020-07-16T13:53:02-04:00
	expiry, err := time.Parse("2006-01-02T15:04:05-07:00", k.String())
	handleError(err)

	if *fErrorOnExpired {
		if time.Until(expiry).Seconds() < 0 {
			os.Exit(1)
		}
	}

	if *fTimeLeft {
		fmt.Printf("%s", time.Until(expiry).Round(time.Second).String())
	} else {
		fmt.Printf("%s", expiry)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Err: %s\n", err)
		os.Exit(1)
	}
}
