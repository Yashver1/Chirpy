package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[string]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	_, err := os.ReadFile(path)
	if errors.Is(err,os.ErrNotExist) {
		if writeErr := os.WriteFile(path, []byte(`{"chirps":{}}`), 0644); writeErr != nil {
			return nil, writeErr
		}
	} else if err != nil {
		return nil, err
	}

	return &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}, nil

}

func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()
	_, err := os.ReadFile(db.path)
	if err == os.ErrNotExist {
		if writeErr := os.WriteFile(db.path, []byte(`{"chirps":{}}`), 0644); writeErr!=nil{
			return writeErr
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.ensureDB()

	db.mux.RLock()
	defer db.mux.RUnlock()

	data, err := os.ReadFile(db.path)
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

	db.mux.Lock()
	defer db.mux.Unlock()

	if err := os.WriteFile(db.path, jsonData, 0644); err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
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
	}

	chirps.Chirps[strings.ToUpper(fmt.Sprintf("%v", count + 1))] = newChirp

	if err := db.writeDB(chirps); err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
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
