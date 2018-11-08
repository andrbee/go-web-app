package models

import (
	"../utils"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)
var (
	UserNotFound        = errors.New("User not found !")
	InvalidLogin        = errors.New("Invalid login !")
	IncorrectAccessData = errors.New("Incorrect data access !")
	UserHasNotCreated   = errors.New("User has not created")
)
type User struct {
	key string
}

func NewUser(username string, hash []byte) (*User, error ){

	id, err := Client.Incr("user:next-id").Result()

	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)

	pipline := Client.Pipeline()
	pipline.HSet(key, "id", id)
	pipline.HSet(key, "username", username)
	pipline.HSet(key, "hash", hash)
	pipline.HSet("user:by-username", username, id)
	pipline.Exec()

	_, err = pipline.Exec()
	if err != nil{
		return nil, UserHasNotCreated
	}

	return &User{key}, nil
}

func (user *User) GetUserName() (string, error){

	return Client.HGet(user.key, "username").Result()
}

func (user *User) GetHash() ([]byte, error){

	return Client.HGet(user.key, "hash").Bytes()
}

func GetUserByUsername(username string) (*User, error) {

	key, err := Client.HGet("user:by-username", username).Int64()
	if err == redis.Nil {
		return nil, UserNotFound
	}
	if err != nil {
		return nil, err
	}
	id := fmt.Sprintf("user:%d", key)

	return &User{id}, nil
}

func AuthUser(username, password string) error{

	user, err := GetUserByUsername(username)
	if err != nil {
		return err
	}

	hashPass, err := user.GetHash()
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashPass, []byte (password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return InvalidLogin
	}

	return err
}

func CreateUser(username, password string) error{

	if username == "" || password == "" {
		return IncorrectAccessData
	}

	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		return utils.SYSTEM_TRY_OPERATION_LATER
	}
	_, err = NewUser(username, hash)

	return err
}