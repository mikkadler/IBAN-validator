package iban

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Country struct {
	Name       string `json:"country_name"`
	Code       string `json:"country_code"`
	IbanLength int    `json:"iban_length"`
	Banks      []Bank
}

type Bank struct {
	Name   string `json:"bank_name"`
	IdCode string `json:"identity_code"`
}

var CountryList []Country

func InitIbanData(folder string) error {
	filename := folder + "countries.json"
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &CountryList)
	if err != nil {
		return err
	}

	for i, country := range CountryList {
		filename = folder + country.Code + "_banks.json"
		data, err := os.ReadFile(filename) // ignore error here since most of the countries are not supported here in this implementaton.
		var bankList []Bank
		if err == nil {
			err = json.Unmarshal(data, &bankList)
			if err != nil {
				fmt.Println(err)
			}
		}
		CountryList[i].Banks = bankList
	}

	return nil
}

func ValidateIBAN(iban string) (bool, error) {
	// Step 0: Let's avoid index errors, check for invalid characters and make all characters uppercase
	if len(iban) < 4 {
		return false, fmt.Errorf("IBAN too short for processing: %s", iban)
	}

	if !ContainsOnlyAlphanumeric(iban) {
		return false, fmt.Errorf("invalid characters present in IBAN string")
	}

	iban = strings.ToUpper(iban)

	// Step 1: Check that the total IBAN length is correct as per the country
	countrySupported := false
	for _, country := range CountryList {
		if iban[:2] == country.Code {
			if country.IbanLength != len(iban) {
				return false, fmt.Errorf("country IBAN length mismatch, is %d, should be %d", len(iban), country.IbanLength)
			}
			countrySupported = true
			break
		}
	}

	if !countrySupported {
		return false, fmt.Errorf("country not supported: %s", iban[:2])
	}

	// Step 2: Move the four initial characters to the end of the string
	iban = iban[4:] + iban[:4]

	// Step 3: Replace each letter in the string with two digits, thereby expanding the string, where A = 10, B = 11, ..., Z = 35
	var digits string
	for _, char := range iban {
		if char >= 'A' && char <= 'Z' {
			digits += strconv.Itoa(int(char - 'A' + 10))
		} else {
			digits += string(char)
		}
	}

	// Step 4: Interpret the string as a decimal big integer and compute the remainder of that number on division by 97
	numericIBAN, ok := new(big.Int).SetString(digits, 10)
	if !ok {
		return false, fmt.Errorf("error converting IBAN to number: %s", digits)
	}
	mod97 := new(big.Int).Mod(numericIBAN, big.NewInt(97))

	if mod97.Cmp(big.NewInt(1)) != 0 {
		return false, fmt.Errorf("checksum failure")
	}
	return true, nil
}

func TryGetBankName(iban string) string {
	//exeptions need to be handled first, currently not supported
	if len(iban) < 4 {
		return "unknown"
	}

	startPosition := 4

	for _, country := range CountryList {
		if iban[:2] == country.Code {
			for _, bank := range country.Banks {
				if len(iban) < startPosition+len(bank.IdCode) {
					return "unknown"
				}
				if iban[startPosition:startPosition+len(bank.IdCode)] == bank.IdCode {
					return bank.Name
				}
			}
		}
	}

	return "unknown"
}

func ContainsOnlyAlphanumeric(input string) bool {
	// Regular expression to match alphanumeric characters
	pattern := "^[a-zA-Z0-9]*$"

	regex := regexp.MustCompile(pattern)

	return regex.MatchString(input)
}
