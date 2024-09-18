package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)


type DataBase interface {
    ensureDB() error
    loadDB() (DBStructure, error)
    writeDB(DBStructure) error
    CreateChirp(string) (Chirp, error)
    GetChirps() ([]Chirp, error)
}

type User struct {
	Id int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
}


type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
	UserId int `json:"author_id"`
}

type Token struct {
	Token string `json:"token"`
	Expires time.Time `json:"expires"`
	UserId int `json:"userid"`
}

type DB struct {
	Path string
	Mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[string]Chirp `json:"chirps"`
	Users map[string]User `json:"users"`
	Tokens map[string]Token `json:"tokens"`
}


func NewDB(path string) (*DB, error) {
	_, err := os.ReadFile(path)
	if errors.Is(err,os.ErrNotExist) {
		if writeErr := os.WriteFile(path, []byte(`{"chirps":{},"users":{},"tokens":{}}`), 0644); writeErr != nil {
			return nil, writeErr
		}
	} else if err != nil {
		return nil, err
	}

	return &DB{
		Path: path,
		Mux:  &sync.RWMutex{},
	}, nil

}

func (db *DB) ensureDB() error {
	db.Mux.Lock()
	defer db.Mux.Unlock()
	_, err := os.ReadFile(db.Path)
	if err == os.ErrNotExist {
		if writeErr := os.WriteFile(db.Path, []byte(`{"chirps":{},"users":{},"tokens":{}}`), 0644); writeErr!=nil{
			return writeErr
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.ensureDB()

	db.Mux.RLock()
	defer db.Mux.RUnlock()

	data, err := os.ReadFile(db.Path)
	if err != nil {
		fmt.Println(1,err)
		return DBStructure{}, err
	}

	var dbstructure DBStructure
	if err := json.Unmarshal(data, &dbstructure); err != nil {
		fmt.Println(2,err)
		return DBStructure{}, err
	}

	return dbstructure, nil

}

func (db *DB) writeDB(dbstructure DBStructure) error {
	db.ensureDB()
	jsonData, err := json.Marshal(dbstructure)
	if err != nil {
		return err
	}

	db.Mux.Lock()
	defer db.Mux.Unlock()

	if err := os.WriteFile(db.Path, jsonData, 0644); err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateUser(user User)(User, error){

	users, err := db.loadDB()
	if err!=nil{
		return User{},err
	}

	max := 0 
	for _,value := range users.Users{
		if value.Id > max{
			max = value.Id
		}
	}

	for _,value := range users.Users{
		if value.Email == user.Email{
			return User{}, fmt.Errorf("User of email: %s exists already",user.Email)
		}
	}

	count := max

	newUser := User{
		Id: count + 1,
		Email: user.Email,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err!=nil{
		return User{}, err
	}

	account := User{
		Id: count + 1,
		Email: user.Email,
		Password: string(hashedPassword),
	}

	users.Users[fmt.Sprintf("%v",count + 1)] = account

	if err := db.writeDB(users); err!=nil{
		return User{}, err
	}

	return newUser,nil
}

func (db *DB) CreateToken(userId int, token string) error{
	tokens, err := db.loadDB()
	if err!=nil{
		return err
	}

	refreshToken := Token{
		Token: token,
		Expires: time.Now().Add(time.Hour * 24 * 60),
		UserId: userId,
	}


	tokens.Tokens[refreshToken.Token] = refreshToken 

	if err:= db.writeDB(tokens); err!=nil{
		return err
	}

	return nil

}

func (db *DB) GetToken(refreshtoken string) (int, error){
	tokens, err := db.loadDB()
	if err!=nil{
		return 0, err
	}

	for _, value := range tokens.Tokens{
		if value.Token == refreshtoken && value.Expires.After(time.Now()){
			return value.UserId,nil
		}
	}

	return 0, fmt.Errorf("Token not found")

}

func (db *DB) DeleteToken(refreshtoken string) error{
	tokens, err := db.loadDB()
	if err!=nil{
		return err
	}

	delete(tokens.Tokens, refreshtoken)
	if err:= db.writeDB(tokens); err!=nil{
		return err
	}
	return nil
}



func (db *DB) UpdateUser(id int, newEmail string, newPassword string) error{
	users , err:= db.loadDB()
	if err!=nil{
		return err
	}

	var foundUser User

	for _,value := range users.Users{
		if value.Id == id{
			foundUser = value
		}

	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword),bcrypt.DefaultCost)
	if err!=nil{
		return err
	}

	users.Users[strconv.Itoa(foundUser.Id)] = User{
		Id: foundUser.Id,
		Email: newEmail ,
		Password: string(hashedPassword),
	}

	if err:= db.writeDB(users); err!=nil{
		return err
	}

	return nil
}

func (db *DB) GetUsers() ([]User, error){

	users, err := db.loadDB()
	if err!=nil{
		return []User{},err
	}

	var userArray []User
	for _,value := range users.Users{
		userArray = append(userArray, value)
	}

	return userArray,nil

}

func (db *DB) CreateChirp(body string, userId int) (Chirp, error) {
	chirps, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	max := 0
	for _,value := range chirps.Chirps{
		if value.Id > max{
			max = value.Id
		}
	}

	count := max

	newChirp := Chirp{
		Id:   count + 1,
		Body: body,
		UserId: userId,
	}

	chirps.Chirps[strings.ToUpper(fmt.Sprintf("%v", count + 1))] = newChirp

	if err := db.writeDB(chirps); err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) DeleteChirp(id string, userId int) error{
	chirps, err := db.loadDB()
	if err!=nil{
		return err
	}

	if chirps.Chirps[id].UserId != userId{
		return fmt.Errorf("User does not have permission to delete this chirp")
	}

	delete(chirps.Chirps,id)

	if err:= db.writeDB(chirps); err!=nil{
		return err
	}

	return nil
	
}

func (db *DB) GetChirps() ([]Chirp, error) {

	chirps, err := db.loadDB()
	if err != nil {
		return []Chirp{}, nil
	}
	var chirpArray []Chirp

	for _, value := range chirps.Chirps {
		chirpArray = append(chirpArray, value)
	}
	return chirpArray, nil
}

func (db *DB) GetChirp(id string) (Chirp, error){
	chirps, err := db.loadDB()
	if err!=nil{
		return Chirp{}, err
	}

	chirp,ok := chirps.Chirps[id]
	if !ok{
		return Chirp{},fmt.Errorf("Chirp with id:%s does not exist",id)
	}

	return chirp, nil

}