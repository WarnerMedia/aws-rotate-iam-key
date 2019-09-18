package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// RTC -- runtime config
type RTC struct {
	Opts map[string]*string
	Flag map[string]*bool
	Vals []string
}

func establishDefaults(rtc *RTC) {
	usr, _ := user.Current()
	credsini := fmt.Sprintf("%s/.aws/credentials", usr.HomeDir)
	rtc.Opts["credsini"] = &credsini
	profnone := "none"
	rtc.Opts["profile"] = &profnone
	output := "none"
	rtc.Opts["output"] = &output
}

func parseCli(rtc *RTC) {
	// take args from command line
	cli := map[string]*string{}
	cli["acct"] = new(string)
	cli["role"] = new(string)

	// get the values from the command line
	cli["key"] = flag.String("k", "", "AWS IAM key.")
	cli["secret"] = flag.String("s", "", "AWS IAM secret")
	cli["credsFile"] = flag.String("c", "", "AWS credentials file")
	cli["profile"] = flag.String("profile", "", "Named profile within AWS credentials file.")
	cli["output"] = flag.String("o", "", "Output format - default is text, option json is json string, /path/to/file runs a regex on the file specified.")
	rtc.Flag["version"] = flag.Bool("v", false, GetVersion())

	flag.Parse()

	appname := filepath.Base(os.Args[0])
	cli["appname"] = &appname

	// overide from cli
	for key := range cli {
		if len(*cli[key]) > 0 {
			rtc.Opts[key] = cli[key]
		}
	}
}

func osenvvar(ev string) *string {
	v := os.Getenv(ev)
	return (&v)
}

// Newrtc - Gets the runtime config together
func Newrtc() *RTC {
	rtc := new(RTC)
	rtc.Opts = make(map[string]*string)
	rtc.Flag = make(map[string]*bool)

	establishDefaults(rtc)

	parseCli(rtc)

	if *rtc.Flag["version"] {
		fmt.Println(GetVersion())
		os.Exit(0)
	}

	return (rtc)
}
