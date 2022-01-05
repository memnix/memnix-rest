package models

// ResponseHTTP structure
type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Count   int         `json:"count"`
}

type CardResponse struct {
	CardID   uint `json:"card_id" example:"1"`
	Card     Card
	Response string `json:"response" example:"42"`
}

type CardResponseValidation struct {
	Validate bool   `json:"validate" example:"true"`
	Message  string `json:"message" example:"Correct answer"`
	Answer   string `json:"correct_answer" example:"42"`
}

type ResponseCard struct {
	Card    Card
	Answers []string
}

type ResponseAuth struct {
	Success bool
	User    User
	Message string
}

type ResponseDeck struct {
	DeckID          uint `json:"deck_id" example:"1"`
	Deck            Deck
	Permission      AccessPermission `json:"permission" example:"1"`
	CardCount       int64            `json:"card_count" example:"42"`
	AverageRating   float32          `json:"average_rating" example:"3.4"`
	PersonnalRating uint             `json:"personnal_rating" example:"3.4"`
	OwnerId         uint             `json:"owner_id" example:"6"`
	Owner           User
}
