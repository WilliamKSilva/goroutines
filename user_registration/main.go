package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type CSVUserField int

const (
	ID CSVUserField = iota
	Name
	Email
	CEP
)

type RawUser struct {
	ID    int
	Name  string
	Email string
	CEP   string
}

type User struct {
	ID      int
	Name    string
	Email   string
	Address struct {
		Street string
		City   string
		CEP    string
	}
}

type AddressViaCepResp struct {
	CEP        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Localidade string `json:"localidade"`
}

func getUsersAddress(users chan User, wg *sync.WaitGroup, rawUsers []RawUser) {
	// If a wait group is passed, this function will be called
	// on a separate Goroutine
	if wg != nil {
		defer wg.Done()
	}

	for _, u := range rawUsers {
		resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", u.CEP))
		if err != nil {
			log.Printf("[ERROR] error tryng to search user address: %v", err)
			continue
		}

		var address AddressViaCepResp
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[ERROR] error trying to read body of address request: %v", err)
			continue
		}

		err = json.Unmarshal(data, &address)
		if err != nil {
			log.Printf("[ERROR] error trying to unmarshal body of address request: %v", err)
			continue
		}

		user := User{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Address: struct {
				Street string
				City   string
				CEP    string
			}{
				Street: address.Logradouro,
				City:   address.Localidade,
				CEP:    address.CEP,
			},
		}

		// If a wait group is passed, this function will be called
		// on a separate Goroutine. So we should use the users channel
		// to send the data to the main thread.
		if wg != nil {
			users <- user
			continue
		}

		log.Printf("[INFO] user without concurrency: %v", user)
	}
}

func getCSVFieldValue(record []string, field CSVUserField) string {
	switch field {
	case ID:
		return record[0]
	case Name:
		return record[1]
	case Email:
		return record[2]
	case CEP:
		return record[3]
	}

	return ""
}

func readUsers(users *[]RawUser, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("[ERROR] error trying to open CSV file: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("[ERROR] error trying to read the data from the CSV file: %v", err)
		return
	}

	for i, record := range records {
		// Ignore the header of the CSV file
		if i == 0 {
			continue
		}

		idStr := getCSVFieldValue(record, ID)
		id, err := strconv.Atoi(getCSVFieldValue(record, ID))
		if err != nil {
			log.Printf("[ERROR] error trying to convert the ID %v of the User: %v", idStr, err)
			continue
		}

		rawUser := RawUser{
			ID:    id,
			Name:  getCSVFieldValue(record, Name),
			Email: getCSVFieldValue(record, Email),
			CEP:   getCSVFieldValue(record, CEP),
		}

		*users = append(*users, rawUser)
	}

}

func main() {
	concurrency := flag.Bool("concurrency", true, "The example should use Goroutines or not")
	threads := flag.Int("threads", 4, "How many Goroutines to spawn")

	flag.Parse()

	t := time.Now()

	var wg sync.WaitGroup
	rawUsers := []RawUser{}
	users := make(chan User)

	readUsers(&rawUsers, "user_registration/users1.csv")
	readUsers(&rawUsers, "user_registration/users2.csv")

	if *concurrency {
		usersPerThread := len(rawUsers) / *threads
		for i := range *threads {
			wg.Add(1)
			go getUsersAddress(users, &wg, rawUsers[i:i+usersPerThread])
		}

		go func() {
			wg.Wait()
			close(users)
		}()

		for user := range users {
			log.Printf("[INFO] user with concurrency: %v", user)
		}
	} else {
		getUsersAddress(users, nil, rawUsers)
	}

	log.Printf("[INFO] time elapsed: %f", time.Since(t).Seconds())
}
