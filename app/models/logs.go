package models

import (
	"encoding/json"
	"fmt"
	"github.com/memnix/memnixrest/pkg/database"
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

func (l *Log) SendLog() error {

	jsonObject, _ := l.ToJson()
	key := fmt.Sprintf("%s.%s", l.Type, l.Event)

	err := database.SendMessageToChannel(database.RabbitMqChan, jsonObject, key)
	return err
}

func (l *Log) Set(Type LogType, Message string, Event LogEvent, UserID, DeckID, CardID uint) {
	l.Type = Type
	l.Message = Message
	l.Event = Event
	l.UserID = UserID
	l.CardID = CardID
	l.DeckID = DeckID
	l.CreatedAt = time.Now()

}

func CreateLog(message string, event LogEvent) *Log {

	return &Log{Message: message, Event: event, CreatedAt: time.Now()}
}

func (l *Log) SetType(Type LogType) *Log {
	l.Type = Type
	return l
}

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

func (l *Log) ToJson() ([]byte, error) {
	body, err := json.Marshal(l)
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
	LogSubscribe           LogEvent = "user.subscribe"
	LogUnsubscribe         LogEvent = "user.unsubscribe"
	LogUserDeckLimit       LogEvent = "user.deckLimit"
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
