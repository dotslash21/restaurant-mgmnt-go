package models

import "github.com/dgrijalva/jwt-go"

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}
