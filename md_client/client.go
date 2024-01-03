package md_client

import (
	"log"

	"github.com/quickfixgo/quickfix"
)

type MarketDataClient struct {
	ApiKey        string
	Passphrase    string
	SecretKey     string
	Subscriptions []string
}

func (mdc MarketDataClient) OnCreate(sessionID quickfix.SessionID) {}

// Upon login subscribe to all symbols
func (mdc MarketDataClient) OnLogon(sessionID quickfix.SessionID) {
	log.Println("logon successful, session id: ", sessionID)

	subscriptionRequest := MarketDataMsgGenerator{}.generate(generateRandomString(5), mdc.Subscriptions)
	if err := quickfix.SendToTarget(subscriptionRequest, sessionID); err != nil {
		log.Println("error while subscribing to symbols, error: ", err)
	}
}

func (mdc MarketDataClient) OnLogout(sessionID quickfix.SessionID) {}

func (mdc MarketDataClient) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	if msg.IsMsgTypeOf("A") {
		log.Println("Sending logon")
		LogonProcessor{}.process(msg, mdc.ApiKey, mdc.Passphrase, mdc.SecretKey)
	}
	log.Println("ToAdmin msg: ", msg)
}

func (mdc MarketDataClient) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error {
	log.Println("ToApp msg: ", msg)
	return nil
}

func (mdc MarketDataClient) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	log.Println("FromAdmin msg: ", msg.ToMessage())
	return nil
}

func (mdc MarketDataClient) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	log.Println("FromApp msg: ", msg)
	return
}
