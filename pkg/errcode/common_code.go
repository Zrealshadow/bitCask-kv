package errcode

var (
	Success                   = NewError(0, "Success")
	ServerError               = NewError(10000000, "Server Error")
	InvalidParams             = NewError(10000001, "Invalid Args")
	NotFound                  = NewError(10000002, "Not Found")
	UnauthorizedAuthNotExist  = NewError(10000003, "Unauthorized Author not Exist")
	UnauthorizedTokenError    = NewError(10000004, "Unauthorized Token Error")
	UnauthorizedTokenTimeout  = NewError(10000005, "Unauthorized Token Timeout")
	UnauthorizedTokenGenerate = NewError(10000006, "Unauthorized Token Generation Failed")
	TooManyRequests           = NewError(10000007, "Too Many Request")
)
