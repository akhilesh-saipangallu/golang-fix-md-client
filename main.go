package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"text/template"

	"github.com/falconxio/fix_md_client/configs"
	"github.com/falconxio/fix_md_client/md_client"
	"github.com/quickfixgo/quickfix"
)

const (
	API_KEY    string = ""
	PASSPHRASE string = ""
	SECRET_KEY string = ""

	FIX_HOST            string = ""
	SENDER_COMP_ID      string = ""
	TARGET_COMP_ID      string = ""
	SOCKET_CONNECT_PORT string = ""
)

var (
	SUBSCRIPTIONS []string = []string{"BTC/USD", "ETH/USD"}
)

func GetConfig(sessionConfig configs.FixSessionConfig) (io.Reader, error) {
	configTemplateFileName := path.Join("configs", "md_client.cfg")
	templateContent, err := ioutil.ReadFile(configTemplateFileName)
	if err != nil {
		return nil, err
	}

	templateString := string(templateContent)
	tmpl, err := template.New("config").Parse(templateString)
	if err != nil {
		return nil, err
	}

	var configString bytes.Buffer
	err = tmpl.Execute(&configString, sessionConfig)
	if err != nil {
		return nil, err
	}
	rawBytes := configString.Bytes()

	return bytes.NewReader(rawBytes), nil
}

func main() {
	config, configErr := GetConfig(configs.FixSessionConfig{
		FixHost:           FIX_HOST,
		FileLogPath:       path.Join("fix_logs"),
		FileStorePath:     path.Join("fix_logs", "store"),
		SenderCompID:      SENDER_COMP_ID,
		TargetCompID:      TARGET_COMP_ID,
		SocketConnectPort: SOCKET_CONNECT_PORT,
	})

	if configErr != nil {
		log.Println("configErr: ", configErr)
	}

	settings, err := quickfix.ParseSettings(config)
	if err != nil {
		log.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	app := md_client.MarketDataClient{
		ApiKey:        API_KEY,
		Passphrase:    PASSPHRASE,
		SecretKey:     SECRET_KEY,
		Subscriptions: SUBSCRIPTIONS,
	}
	fileLogFactory, err := quickfix.NewFileLogFactory(settings)
	if err != nil {
		log.Printf("Error creating FileLogFactory: %s\n", err)
		os.Exit(1)
	}

	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), settings, fileLogFactory)
	if err != nil {
		log.Printf("Unable to create Initiator: %s\n", err)
		os.Exit(1)
	}
	initiator.Start()
	defer initiator.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
