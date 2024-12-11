package md_client

import (
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	fix44mdr "github.com/quickfixgo/fix44/marketdatarequest"
)

type MarketDataSubMsgGenerator struct{}

func (MarketDataSubMsgGenerator) generate(mdReqID string, symbols []string) fix44mdr.MarketDataRequest {
	request := fix44mdr.New(
		field.NewMDReqID(mdReqID),
		field.NewSubscriptionRequestType(enum.SubscriptionRequestType_SNAPSHOT_PLUS_UPDATES),
		field.NewMarketDepth(0),
	)

	entryTypes := fix44mdr.NewNoMDEntryTypesRepeatingGroup()
	entryTypes.Add().SetMDEntryType(enum.MDEntryType_BID)
	entryTypes.Add().SetMDEntryType(enum.MDEntryType_OFFER)
	request.SetNoMDEntryTypes(entryTypes)

	relatedSym := fix44mdr.NewNoRelatedSymRepeatingGroup()
	for _, symbol := range symbols {
		relatedSym.Add().SetSymbol(symbol)
	}
	request.SetNoRelatedSym(relatedSym)
	request.SetMDUpdateType(enum.MDUpdateType_INCREMENTAL_REFRESH)

	return request
}

type MarketDataUnSubMsgGenerator struct{}

func (MarketDataUnSubMsgGenerator) generate(mdReqID string, symbols []string) fix44mdr.MarketDataRequest {
	request := fix44mdr.New(
		field.NewMDReqID(mdReqID),
		field.NewSubscriptionRequestType(enum.SubscriptionRequestType_DISABLE_PREVIOUS_SNAPSHOT_PLUS_UPDATE_REQUEST),
		field.NewMarketDepth(0),
	)

	entryTypes := fix44mdr.NewNoMDEntryTypesRepeatingGroup()
	entryTypes.Add().SetMDEntryType(enum.MDEntryType_BID)
	entryTypes.Add().SetMDEntryType(enum.MDEntryType_OFFER)
	request.SetNoMDEntryTypes(entryTypes)

	relatedSym := fix44mdr.NewNoRelatedSymRepeatingGroup()
	for _, symbol := range symbols {
		relatedSym.Add().SetSymbol(symbol)
	}
	request.SetNoRelatedSym(relatedSym)
	request.SetMDUpdateType(enum.MDUpdateType_INCREMENTAL_REFRESH)

	return request
}
