package security

import "golang.org/x/crypto/bcrypt"

func Hash(senha string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

func VerificarSenha(senhaHash string, senha string) error {
	return bcrypt.CompareHashAndPassword([]byte(senhaHash), []byte(senha))
}
