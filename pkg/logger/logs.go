package logger

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/memnix/memnixrest/data/infrastructures"
	"time"
)

type Log struct {
	Type      LogType   `json:"type"`
	Message   string    `json:"message"`
	Event     LogEvent  `json:"event"`
	CreatedAt time.Time `json:"createdat"`
	UserID    uint      `json:"userid"`
	DeckID    uint      `json:"deckid"`
	CardID    uint      `json:"cardid"`
}

// SendLog sends log to the rabbitmq channel
func (l *Log) SendLog() error {
	jsonObject, _ := l.ToJSON()
	key := fmt.Sprintf("%s.%s", l.Type, l.Event)

	err := infrastructures.SendMessageToChannel(infrastructures.RabbitMQ.Channel, jsonObject, key)
	return err
}

// Set Log
func (l *Log) Set(logType LogType, message string, event LogEvent, userID, deckID, cardID uint) {
	l.Type = logType
	l.Message = message
	l.Event = event
	l.UserID = userID
	l.CardID = cardID
	l.DeckID = deckID
	l.CreatedAt = time.Now()
}

// CreateLog returns a new Log object
func CreateLog(message string, event LogEvent) *Log {
	return &Log{Message: message, Event: event, CreatedAt: time.Now()}
}

// SetType method
func (l *Log) SetType(logType LogType) *Log {
	l.Type = logType
	return l
}

// AttachIDs method
func (l *Log) AttachIDs(userID, deckID, cardID uint) *Log {
	if userID != 0 {
		l.UserID = userID
	}

	if cardID != 0 {
		l.CardID = cardID
	}
	if deckID != 0 {
		l.DeckID = deckID
	}

	return l
}

// ToJSON method
func (l *Log) ToJSON() ([]byte, error) {
	body, err := sonic.Marshal(l)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type LogType string

const (
	LogTypeInfo    LogType = "info"
	LogTypeWarning LogType = "warning"
	LogTypeError   LogType = "error"
)

// LogEvent enum type
type LogEvent string

const (
	LogUndefined           LogEvent = "undefined"
	LogUserLogin           LogEvent = "user.login"
	LogUserLogout          LogEvent = "user.logout"
	LogUserRegister        LogEvent = "user.register"
	LogUserEdit            LogEvent = "user.edit"
	LogUserDeleted         LogEvent = "user.deleted"
	LogUserPasswordReset   LogEvent = "user.password_reset"
	LogUserPasswordChanged LogEvent = "user.password_changed"
	LogSubscribe           LogEvent = "user.subscribe"
	LogUnsubscribe         LogEvent = "user.unsubscribe"
	LogUserDeckLimit       LogEvent = "user.deckLimit"
	LogPublishRequest      LogEvent = "deck.publish"
	LogDeckCreated         LogEvent = "deck.created"
	LogDeckDeleted         LogEvent = "deck.deleted"
	LogDeckEdited          LogEvent = "deck.edited"
	LogDeckCardLimit       LogEvent = "deck.cardLimit"
	LogCardCreated         LogEvent = "card.created"
	LogCardDeleted         LogEvent = "card.deleted"
	LogCardEdited          LogEvent = "card.edited"
	LogAlreadyUsedEmail    LogEvent = "register.usedEmail"
	LogIncorrectEmail      LogEvent = "login.incorrectEmail"
	LogIncorrectPassword   LogEvent = "login.incorrectPassword"
	LogLoginError          LogEvent = "login.error"
	LogPermissionForbidden LogEvent = "permission.forbidden"
	LogQueryGetError       LogEvent = "query.get"
	LogBodyParserError     LogEvent = "query.bodyParser"
	LogBadRequest          LogEvent = "query.badRequest"
)
