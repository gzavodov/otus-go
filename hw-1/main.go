package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

// NtpHostURL URL of active NTP host
const NtpHostURL = "ru.pool.ntp.org"

func main() {
	time, err := ntp.Time(NtpHostURL)
	if err != nil {
		//This code may be replaced by log.Fatal
		fmt.Fprintf(os.Stderr, "ERROR: '%s'", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Date: %s\n", time.Format("02 January 2006"))
	fmt.Fprintf(os.Stdout, "Time: %s\n", time.Format("15:04:05 MST"))
}
