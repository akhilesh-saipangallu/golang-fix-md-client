package md_client

import (
	"log"

	"github.com/quickfixgo/quickfix"
)

type MarketDataClient struct {
	ApiKey          string
	Passphrase      string
	SecretKey       string
	SessionId       quickfix.SessionID
	Subscriptions   []string
	SubscriptionIds []string
	Routes          *quickfix.MessageRouter
}

func (mdc *MarketDataClient) OnCreate(sessionID quickfix.SessionID) {
	mdc.SessionId = sessionID
}

// Upon login subscribe to all symbols
func (mdc *MarketDataClient) OnLogon(sessionID quickfix.SessionID) {
	log.Println("logon successful, session id: ", sessionID)

	subReqId := generateRandomString(5)
	subscriptionRequest := MarketDataSubMsgGenerator{}.generate(subReqId, mdc.Subscriptions)
	if err := quickfix.SendToTarget(subscriptionRequest, sessionID); err != nil {
		log.Println("error while subscribing symbols, error: ", err)
	}

	// track the reqId
	mdc.SubscriptionIds = append(mdc.SubscriptionIds, subReqId)
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

func (mdc *MarketDataClient) UnSubscribeAll() {
	for _, subReqId := range mdc.SubscriptionIds {
		subscriptionRequest := MarketDataUnSubMsgGenerator{}.generate(subReqId, mdc.Subscriptions)

		if err := quickfix.SendToTarget(subscriptionRequest, mdc.SessionId); err != nil {
			log.Println("error while unsubscribing symbols, error: ", err)
		}
	}
}
