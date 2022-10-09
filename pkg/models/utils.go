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
	CardID   uint   `json:"card_id" example:"1"`
	Card     Card   `json:"-" swaggerignore:"true"`
	Response string `json:"response" example:"42"`
	Training bool   `json:"training" example:"false"`
}

type CardSelfResponse struct {
	Training bool `json:"training" example:"false"`
	Quality  uint `json:"quality" example:"3"` // Min 0 - Max 4
	CardID   uint `json:"card_id" example:"1"`
	Card     Card
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
	Card          Card
	Answers       []string
	LearningStage LearningStage `json:"learning_stage"`
}

// Set ResponseCard values
func (responseCard *ResponseCard) Set(memdate *MemDate, answers []string) {
	responseCard.Answers = answers
	responseCard.Card = memdate.Card
	responseCard.LearningStage = memdate.LearningStage
}

// ResponseAuth struct
type ResponseAuth struct {
	Success bool   `json:"success"`
	User    User   `json:"user"`
	Message string `json:"message"`
}

type ResponsePublicAuth struct {
	Success bool       `json:"success"`
	User    PublicUser `json:"user"`
	Message string     `json:"message"`
}

// ResponseDeck struct
type ResponseDeck struct {
	DeckID      uint             `json:"deck_id" example:"1"`
	Permission  AccessPermission `json:"permission" example:"1"`
	CardCount   uint16           `json:"card_count" example:"42"`
	ToggleToday bool             `json:"settings_today" `
	OwnerID     uint             `json:"owner_id" example:"6"`
	Owner       PublicUser       `swaggerignore:"true" json:"Owner,omitempty"`
	Deck        Deck             `json:"Deck,omitempty"`
}
