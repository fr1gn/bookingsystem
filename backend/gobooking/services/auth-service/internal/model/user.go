package model

type User struct {
	ID            string `bson:"_id,omitempty"`
	FullName      string `bson:"full_name"`
	Email         string `bson:"email"`
	PasswordHash  string `bson:"password_hash"`
	EmailVerified bool   `bson:"email_verified"`
}
