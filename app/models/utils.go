package models

// ResponseHTTP structure to format API answers
type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Count   int         `json:"count"`
}

// Set ResponseHTTP values
func (res *ResponseHTTP) Set(success bool, message string, data interface{}, count int) {
	*res = ResponseHTTP{
		Success: success,
		Data:    data,
		Message: message,
		Count:   count,
	}
}

// GenerateError method
func (res *ResponseHTTP) GenerateError(message string) {
	res.Set(false, message, nil, 0)
}

// GenerateSuccess method
func (res *ResponseHTTP) GenerateSuccess(message string, data interface{}, count int) {
	res.Set(true, message, data, count)
}

// DeckConfig struct
type DeckConfig struct {
	TodaySetting bool `json:"settings_today"`
}

// CardResponse struct
type CardResponse struct {
	CardID   uint `json:"card_id" example:"1"`
	Card     Card
	Response string `json:"response" example:"42"`
	Training bool   `json:"training" example:"false"`
}

// CardResponseValidation struct
type CardResponseValidation struct {
	Validate bool   `json:"validate" example:"true"`
	Message  string `json:"message" example:"Correct answer"`
	Answer   string `json:"correct_answer" example:"42"`
}

func (validation *CardResponseValidation) SetCorrect() {
	validation.Validate = true
	validation.Message = "Correct answer"
}

func (validation *CardResponseValidation) SetIncorrect() {
	validation.Validate = false
	validation.Message = "Incorrect answer"
}

// ResponseCard struct
type ResponseCard struct {
	Card    Card
	Answers []string
}

// Set ResponseCard values
func (responseCard *ResponseCard) Set(card *Card, answers []string) {
	responseCard.Answers = answers
	responseCard.Card = *card
}

// ResponseAuth struct
type ResponseAuth struct {
	Success bool
	User    User
	Message string
}

// ResponseDeck struct
type ResponseDeck struct {
	DeckID      uint `json:"deck_id" example:"1"`
	Deck        Deck
	Permission  AccessPermission `json:"permission" example:"1"`
	CardCount   int64            `json:"card_count" example:"42"`
	OwnerId     uint             `json:"owner_id" example:"6"`
	Owner       PublicUser
	ToggleToday bool `json:"settings_today" `
}
