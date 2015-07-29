package main

import (
	"database/sql"
	_ "expvar"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/bugsnag/bugsnag-go"
	_ "github.com/docker/distribution/registry/auth/htpasswd"
	_ "github.com/docker/distribution/registry/auth/token"
	"github.com/endophage/gotuf/signed"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"

	bugsnag_hook "github.com/Sirupsen/logrus/hooks/bugsnag"
	"github.com/docker/notary/server"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/signer"
	"github.com/spf13/viper"
)

// DebugAddress is the debug server address to listen on
const DebugAddress = "localhost:8080"

var debug bool
var configFile string

func init() {
	// set default log level to Error
	viper.SetDefault("logging", map[string]interface{}{"level": 2})

	// Setup flags
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.BoolVar(&debug, "debug", false, "Enable the debugging server on localhost:8080")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if debug {
		go debugServer(DebugAddress)
	}

	ctx := context.Background()

	filename := filepath.Base(configFile)
	ext := filepath.Ext(configFile)
	configPath := filepath.Dir(configFile)

	viper.SetConfigType(strings.TrimPrefix(ext, "."))
	viper.SetConfigName(strings.TrimSuffix(filename, ext))
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error("Viper Error: ", err.Error())
		logrus.Error("Could not read config at ", configFile)
		os.Exit(1)
	}
	logrus.SetLevel(logrus.Level(viper.GetInt("logging.level")))

	// set up bugsnag and attach to logrus
	bugs := viper.GetString("reporting.bugsnag")
	if bugs != "" {
		apiKey := viper.GetString("reporting.bugsnag_api_key")
		releaseStage := viper.GetString("reporting.bugsnag_release_stage")
		bugsnag.Configure(bugsnag.Configuration{
			APIKey:       apiKey,
			ReleaseStage: releaseStage,
		})
		hook, err := bugsnag_hook.NewBugsnagHook()
		if err != nil {
			logrus.Error("Could not attach bugsnag to logrus: ", err.Error())
		} else {
			logrus.AddHook(hook)
		}
	}
	keyAlgo := viper.GetString("trust_service.key_algorithm")
	if keyAlgo == "" {
		logrus.Fatal("no key algorithm configured.")
		os.Exit(1)
	}
	ctx = context.WithValue(ctx, "keyAlgorithm", keyAlgo)

	var trust signed.CryptoService
	if viper.GetString("trust_service.type") == "remote" {
		logrus.Info("[Notary Server] : Using remote signing service")
		trust = signer.NewNotarySigner(
			viper.GetString("trust_service.hostname"),
			viper.GetString("trust_service.port"),
			viper.GetString("trust_service.tls_ca_file"),
		)
	} else {
		logrus.Info("[Notary Server] : Using local signing service")
		trust = signed.NewEd25519()
	}

	if viper.GetString("storage.backend") == "mysql" {
		logrus.Debug("Using mysql backend")
		dbURL := viper.GetString("storage.db_url")
		db, err := sql.Open("mysql", dbURL)
		if err != nil {
			logrus.Fatal("[Notary Server] Error starting DB driver: ", err.Error())
			return // not strictly needed but let's be explicit
		}
		ctx = context.WithValue(ctx, "metaStore", storage.NewMySQLStorage(db))
	} else {
		logrus.Debug("Using memory backend")
		ctx = context.WithValue(ctx, "metaStore", storage.NewMemStorage())
	}
	logrus.Info("[Notary Server] Starting Server")
	err = server.Run(
		ctx,
		viper.GetString("server.addr"),
		viper.GetString("server.tls_cert_file"),
		viper.GetString("server.tls_key_file"),
		trust,
		viper.GetString("auth.type"),
		viper.Get("auth.options"),
	)

	logrus.Error("[Notary Server]", err.Error())
	return
}

func usage() {
	fmt.Println("usage:", os.Args[0])
	flag.PrintDefaults()
}

// debugServer starts the debug server with pprof, expvar among other
// endpoints. The addr should not be exposed externally. For most of these to
// work, tls cannot be enabled on the endpoint, so it is generally separate.
func debugServer(addr string) {
	logrus.Info("[Notary Debug Server] server listening on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logrus.Fatal("[Notary Debug Server] error listening on debug interface: ", err)
	}
}
