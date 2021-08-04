package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/jakewarren/rssfinder"
	"github.com/spf13/pflag"
)

var (
	// build information set by ldflags
	appName    = "rssfinder"
	appVersion = "(ﾉ☉ヮ⚆)ﾉ ⌒*:･ﾟ✧"
	commit     = "(ﾉ☉ヮ⚆)ﾉ ⌒*:･ﾟ✧"
	buildDate  = "(ﾉ☉ヮ⚆)ﾉ ⌒*:･ﾟ✧"
)

const usageMessage = `Usage: rssfinder [flags] <URL>`

func main() {
	r := rssfinder.New()

	pflag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, usageMessage)
		_, _ = fmt.Fprintln(os.Stderr, "")
		_, _ = fmt.Fprintln(os.Stderr, "Flags:")
		pflag.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr, "")
		_, _ = fmt.Fprintln(os.Stderr, "URL: https://github.com/jakewarren/rssfinder")
	}

	displayHelp := pflag.BoolP("help", "h", false, "display help")
	displayVersion := pflag.BoolP("version", "V", false, "display version information")

	pflag.BoolVarP(&r.Config.EnableFuzzer, "fuzzer", "f", false, "enables the fuzzer module")
	pflag.BoolVarP(&r.Config.EnableScraper, "scraper", "s", false, "enables the scraper module")
	pflag.BoolVar(&r.Config.EnableNewsBlur, "newsblur", false, "enables the newblur module")
	pflag.BoolVar(&r.Config.EnableInoreader, "inoreader", true, "enables the inoreader module")
	pflag.BoolVar(&r.Config.EnableFeedly, "feedly", true, "enables the feedly module")
	verbose := pflag.BoolP("verbose", "v", false, "enable trace logging")
	pflag.Parse()

	// override the default usage display
	if *displayHelp {
		pflag.Usage()
		os.Exit(0)
	}

	if *displayVersion {
		fmt.Printf(`%s:
    version     : %s
    git hash    : %s
    build date  : %s 
    go version  : %s
    go compiler : %s
    platform    : %s/%s
`, appName, appVersion, commit, buildDate, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if pflag.NArg() == 0 {
		log.Fatalln("URL not provided")
	}

	if *verbose {
		rssfinder.LogFunc = log.Printf
	}

	feeds := r.FindRSS(pflag.Arg(0))
	jsonOutput(feeds)
}

func jsonOutput(feeds []rssfinder.Feed) {
	out, _ := json.MarshalIndent(feeds, "", "  ")
	fmt.Println(string(out))
}
