package authdb

import (
	"golang.org/x/crypto/bcrypt"
)

// would be a database table in a real application
// And NEVER STORE PASSWORDS IN PLAINTEXT; some kind of hash should always be used
var usersPasswords = map[string][]byte{
	"shien": []byte("$2a$12$aMfFQpGSiPiYkekov7LOsu63pZFaWzmlfm1T8lvG6JFj2Bh4SZPWS"),
	"john":  []byte("$2a$12$l398tX477zeEBP6Se0mAv.ZLR8.LZZehuDgbtw2yoQeMjIyCNCsRW"),
}

func VerifyUserPassword(username string, password string) bool {
	targetPassword, hasPassword := usersPasswords[username]

	if !hasPassword {
		return false
	}

	if cmpErr := bcrypt.CompareHashAndPassword(targetPassword, []byte(password)); cmpErr == nil {
		return true
	}

	return false
}
