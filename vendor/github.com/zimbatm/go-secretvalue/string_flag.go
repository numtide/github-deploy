package secretvalue

// StringFlag implements the flag.Flag interface and can be used during
// command-line option parsing. Can also be used with github.com/urfave/cli
type StringFlag struct {
	*SecretValue
}

// Set records the value as a secret value
func (s StringFlag) Set(value string) error {
	if s.SecretValue == nil {
		s.SecretValue = New("unknonw")
	}
	s.SecretValue.Set([]byte(value))
	return nil
}

func (s StringFlag) String() string {
	return s.SecretValue.String()
}
