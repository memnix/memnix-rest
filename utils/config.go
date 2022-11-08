package utils

const (
	// SQLMaxOpenConns is the maximum number of open connections to the infrastructures.
	SQLMaxOpenConns = 50
	// SQLMaxIdleConns is the maximum number of connections in the idle connection pool.
	SQLMaxIdleConns = 10
	// CacheExpireMinutes is the cache expire minutes.
	CacheExpireMinutes = 2

	MaxDeckNormalUser = 5
	MaxCardDeck       = 200
	MaxMcqDeck        = 100

	MaxPasswordLen = 50
	MaxEmailLen    = 100
	MaxDefaultLen  = 200
	MaxUsernameLen = 50

	MinCardQuestionLen    = 1
	MaxCardFormatLen      = 50
	MaxImageURLLen        = 200
	MaxCardExplicationLen = 500

	MaxDeckNameLen = 42
	MinDeckNameLen = 5
	DeckKeyLen     = 4
	MaxLangLen     = 2

	MaxMcqAnswersLen = 150
	MinMcqAnswersLen = 4
	MaxMcqName       = 50
)
