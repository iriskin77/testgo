package constants

import "time"

type userCtx string
type sortCtx string

const (

	// for AuthMiddleware
	AuthorizationHeader         = "Authorization"
	UserContextKey      userCtx = "userId"

	// for SortMiddleware
	AscSort                   = "ASC"
	DescSort                  = "DESC"
	OptionsContextKey sortCtx = "sort_options"
	DefaulSortField           = "created_at"
	DefaultSortOrder          = "ASC"

	// for token
	TimeTokenExpire = 12 * time.Hour
	SigningKey      = "qrkjk#4#%35FSFJlja#4353KSFjH"

	// salt to hash password
	SaltPassword = "hjqrhjqw124617ajfhajs"

	Host = "http://localhost:8000"
)
