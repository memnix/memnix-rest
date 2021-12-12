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

