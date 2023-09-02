package config

const (
	MB = 1024 * 1024 // MB is the number of bytes in a megabyte

	Base10  = 10 // Base10 is the base 10
	BitSize = 32 // BitSize is the bit size

	JwtTokenHeaderLen = 2 // JwtTokenHeaderLen is the jwt token header length

	MaxPasswordLength = 72 // MaxPasswordLength is the max password length
	MinPasswordLength = 8  // MinPasswordLength is the min password length

	Localhost = "http://localhost:1815"   // Localhost is the localhost url
	APIHost   = "https://beta.memnix.app" // APIHost is the api host url

	FrontHost      = "https://memnix.corentings.dev" // FrontHost is the front host url
	FrontHostLocal = "http://localhost:3000"         // FrontHostLocal is the front host url for local development
)
