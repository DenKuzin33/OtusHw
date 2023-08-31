package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (users, error) {
	scanner := bufio.NewScanner(r)
	i := 0
	var user User
	var result users
	for scanner.Scan() {
		if err := user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return result, err
		}
		result[i] = user
		i++
	}
	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainRegex, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		matched := domainRegex.Match([]byte(user.Email))

		if matched {
			result[strings.ToLower(user.Email[strings.LastIndex(user.Email, "@")+1:])]++
		}
	}
	return result, nil
}
