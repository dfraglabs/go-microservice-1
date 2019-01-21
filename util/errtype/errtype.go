package errtype

//Base error struct
type TypedErr struct {
	E error
}

func (t TypedErr) Error() string { return t.E.Error() }

//Typed errors
type APINotFoundErr TypedErr
type BackendRequestFailed TypedErr
type SNSPublishErr TypedErr
type TokenSigningErr TypedErr

//DB
type DuplicateKeyErr TypedErr
type KeyNotFoundErr TypedErr

//Auth
type InvalidPasswordErr TypedErr
type MissingCredentialsErr TypedErr
type JWTExpiredError TypedErr

//User
type UserNotFoundErr TypedErr
type UserNotConfiguredErr TypedErr
type UserInvalidErr TypedErr

//Session
type SessionNotFoundErr TypedErr
type DomainNotFoundErr TypedErr
type SessionCreationErr TypedErr
type DuplicateSessionErr TypedErr

//Error funcs
func (e APINotFoundErr) Error() string        { return e.E.Error() }
func (e BackendRequestFailed) Error() string  { return e.E.Error() }
func (e SNSPublishErr) Error() string         { return e.E.Error() }
func (e TokenSigningErr) Error() string       { return e.E.Error() }
func (e DuplicateKeyErr) Error() string       { return e.E.Error() }
func (e KeyNotFoundErr) Error() string        { return e.E.Error() }
func (e InvalidPasswordErr) Error() string    { return e.E.Error() }
func (e MissingCredentialsErr) Error() string { return e.E.Error() }
func (e UserNotFoundErr) Error() string       { return e.E.Error() }
func (e UserNotConfiguredErr) Error() string  { return e.E.Error() }
func (e UserInvalidErr) Error() string        { return e.E.Error() }
func (e SessionNotFoundErr) Error() string    { return e.E.Error() }
func (e DomainNotFoundErr) Error() string     { return e.E.Error() }
func (e SessionCreationErr) Error() string    { return e.E.Error() }
func (e DuplicateSessionErr) Error() string   { return e.E.Error() }
func (e JWTExpiredError) Error() string       { return e.E.Error() }
