package types

import "fmt"

var _ fmt.Stringer = (*SecretString)(nil)

type SecretString string

func (s SecretString) String() string {
	return "*****"
}

func (s SecretString) Raw() string {
	return string(s)
}
