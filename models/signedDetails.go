package models

import "github.com/dgrijalva/jwt-go"

type SignedDetails struct {
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Uid       string `json:"uid,omitempty"`
	jwt.StandardClaims
}
