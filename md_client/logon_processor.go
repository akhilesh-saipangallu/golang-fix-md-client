package md_client

import (
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
)

type LogonProcessor struct{}

func (LogonProcessor) process(logonMsg *quickfix.Message, apiKey string, passphrase string, secretKey string) {
	encryptMethod := field.NewEncryptMethod(enum.EncryptMethod_NONE_OTHER)
	logonMsg.Body.Set(encryptMethod)

	usernameField := field.NewUsername(apiKey)
	logonMsg.Body.Set(usernameField)

	passwordField := field.NewPassword(passphrase)
	logonMsg.Body.Set(passwordField)

	sendingTime, _ := logonMsg.Header.GetString(52)
	messageType, _ := logonMsg.Header.GetString(35)
	messageSeqNumber, _ := logonMsg.Header.GetString(34)
	senderCompId, _ := logonMsg.Header.GetString(49)
	targetCompId, _ := logonMsg.Header.GetString(56)
	password, _ := logonMsg.Body.GetString(554)
	rawData := generateSignature(
		secretKey, sendingTime, messageType, messageSeqNumber, senderCompId, targetCompId, password,
	)
	rawDataField := field.NewRawData(rawData)
	logonMsg.Body.Set(rawDataField)

	resetOnLogon := field.NewResetSeqNumFlag(true)
	logonMsg.Body.Set(resetOnLogon)
}
