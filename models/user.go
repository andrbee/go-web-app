package models

import (
	"../utils"
	"errors"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)
var (
	USER_NOT_FOUND = errors.New("User not found !")
	INVALID_LOGIN = errors.New("Invalid login !")
	INCORRECT_ACCESS_DATA = errors.New("Incorrect data access !")
)
func AuthUser(name, password string) error{
	hashPass, err := Client.Get("user:" + name).Bytes()
	if err == redis.Nil {
		return USER_NOT_FOUND
	}

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashPass, []byte (password))

	if err != nil {
		return INVALID_LOGIN
	}

	return nil;

}

func CreateUser(username, password string) error{

	if username == "" || password == "" {
		return INCORRECT_ACCESS_DATA
	}

	cost := bcrypt.DefaultCost
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		return utils.SYSTEM_TRY_OPERATION_LATER
	}

	return Client.Set("user:" + username, hashPass, 0).Err()
}