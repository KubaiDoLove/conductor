package apiserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/KubaiDoLove/conductor/internal/app/database"
	"github.com/KubaiDoLove/conductor/internal/app/database/drivers"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/logutils"
)

// Start is a command to start new api server.
func Start(version string) {
	fmt.Printf("conductor %s\n", version)

	opts := ConfigWithParsedFlags()
	setupLog(opts.InDebugMode)

	appCtx, cancelAppCtx := context.WithCancel(context.Background())
	defer cancelAppCtx()

	// ловим сигнал для graceful termination
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Print("[WARN] interrupt signal")
		cancelAppCtx()
	}()

	ds, err := setupDS(*opts)
	if err != nil {
		log.Println(err)
		return
	}
	defer ds.Close(context.Background())

	serverApp := NewHTTPServer(appCtx, opts, ds, version)
	serverApp.Run()

	log.Printf("[INFO] process terminated")
}

// setupLog sets up log levels and logs output.
func setupLog(inDebugMode bool) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stdout,
	}

	log.SetFlags(log.Ldate | log.Ltime)

	if inDebugMode {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
		filter.MinLevel = "DEBUG"
	}

	log.SetOutput(filter)
}

func setupDS(opts Config) (drivers.DataStore, error) {
	ds, err := database.Connect(drivers.DataStoreConfig{
		URL:           opts.DSURL,
		DataStoreName: opts.DSName,
		DataBaseName:  opts.DSDB,
	})
	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] cannot connect to datastore %s: %v", opts.DSName, err)
		return nil, errors.New(errMsg)
	}

	log.Printf("[INFO] connected to %s", ds.Name())

	return ds, nil
}
