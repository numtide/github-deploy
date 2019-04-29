// Package secretvalue provides a value wrapper to hold secrets. It's main
// purpose is to help avoid sending secrets to logs by mistake.
package secretvalue

// SecretValue is a value holder that avoids accidentally leaking the value
// into the logs
type SecretValue struct {
	name string
	secret []byte
}

// New creates a new SecretValue with the type set
func New(secretName string) *SecretValue {
	return &SecretValue{
		name: secretName,
	}
}

// NewWithValue is used with the immutable style
func NewWithValue(name string, secret []byte) *SecretValue {
	return &SecretValue {
		name: name,
		secret: secret,
	}
}

// String returns only the secret type to avoid logging the secret
func (s SecretValue) String() string {
	return "<secret:" + s.name + ">"
}

// Get returns the secret. Be careful with handling the returned data!
func (s SecretValue) Get() []byte {
	return s.secret
}

// GetString returns the secret as a string. Be careful with handling the
// returned data!
func (s SecretValue) GetString() string {
	return string(s.secret)
}

// Set mutates the SecretValue and inserts the secret into the value holder
func (s *SecretValue) Set(secret []byte) *SecretValue {
	s.secret = secret
	return s
}

// SetString works like Set but also converts the string into the underlying
// []byte storage
func (s *SecretValue) SetString(secret string) *SecretValue {
	s.secret = []byte(secret)
	return s
}

// Forget will destroy the secret, can be used with `defer` to make sure that
// a secret is forgotten when existing a function
func (s *SecretValue) Forget() {
	s.secret = nil
}
