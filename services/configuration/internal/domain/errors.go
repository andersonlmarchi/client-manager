package domain

import "github.com/andersonlmarchi/client-manager/packages/shared"

func errInvalid(msg string) error {
	return shared.NewError(shared.CodeInvalid, msg)
}
