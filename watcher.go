/*
A package comment
*/
package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log" // single, global logger
)

// MyWatcher is merely a convenient alias
// N.B -- capitalized, so exported.
type MyWatcher fsnotify.Watcher // a semi-convenient alias

func main() {

	verbose := flag.Bool("v", false, "enable for more detail")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
		// and exit
	}

	os.Exit(Main(flag.Arg(0), *verbose))
}

// Main Command line monitors a specified directory and dumps
// fsnotify.Write events.
func Main(path string, verbosity bool) int {

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr, TimeFormat: time.RFC3339})

	if verbosity {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	watchPath, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		log.Fatal(). // will exit
				Err(err).
				Str("path", flag.Arg(0)).
				Msg("FilePath error")
		//log.Fatalf("Error on %s: %s", flag.Arg(0), err)
	}

	log.Debug().
		Str("path", watchPath).
		Msg("Watching writes")

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal().Err(err).Msg("NewWatcher err")
	}

	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Info().Str("path", event.Name).Msg("modified file")
				}
			case err := <-watcher.Errors:
				log.Warn().Err(err).Msg("error from watcher")
			}
		}
	}()

	err = watcher.Add(watchPath)
	if err != nil {
		log.Fatal().Err(err).Msg("watcher.add")
	}
	<-done

	return 0
}
