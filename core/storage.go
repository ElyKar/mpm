package core

// Core contains the logic and functionalities to be called by the different user interfaces

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"golang.org/x/crypto/bcrypt"
)

// Storage is used to keep track of both the user master file
type Storage struct {
	// Passphrase of the user. That's a bcrypt hash.
	Passphrase string `json:"Pass"`
	// The different sections of the file. The format is the following:
	//     - section_name_1:
	//         * pass_name_1: encrypted_pass_1
	//         * pass_name_2: encrypted_pass_2
	//     - section_name_2:
	//         * pass_name_1: encrypted_pass_1
	//         * pass_name_2: encrypted_pass_2
	Sections map[string]map[string]string `json:"Sections"`
}

// Path to the master file
var fileName string = path.Join(os.Getenv("HOME"), ".mpm")

// Default bcrypt cost
const bcryptCost = 10

// InitPassphrase creates a new storage with given passphrase
func InitPassphrase(new []byte) *Storage {
	hashed, _ := bcrypt.GenerateFromPassword(new, bcryptCost)

	return &Storage{string(hashed), make(map[string]map[string]string)}
}

// GetStorage reads the master file, and creates the matching storage object if possible. If not, an error is raised (invalid permissions, non-existent file, wrong formatting, etc ...
func GetStorage() (*Storage, error) {
	_, err := os.Stat(fileName)
	if err != nil {
		return nil, fmt.Errorf("File %s does not exists", fileName)
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	store := &Storage{}
	err = json.Unmarshal(data, store)

	if err != nil {
		return nil, err
	}

	return store, nil
}

// Saves the file on the disk, raises an error if necessary with JSON formatting
func (s *Storage) DumpOnDisk() error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, data, 0600)
	return err
}

// CheckPassphrase verifies the given passphrase against the stored hash using bcrypt's hash function.
func (s *Storage) CheckPassphrase(pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(s.Passphrase), []byte(pass))
}

// SetNewPassphrase changes from the old passphrase to the new one. It checks for validity of the old one first, then decrypts all passwords and re-encrypts them with the new passphrase. If there is an error during this operation, nothing is comitted.
func (s *Storage) SetNewPassphrase(old string, new string) error {
	err := bcrypt.CompareHashAndPassword([]byte(s.Passphrase), []byte(old))
	if err != nil {
		return err
	}

	// Decode with old one, encode with new one.
	dec, enc := NewTranscoder(old), NewTranscoder(new)

	// Iterate through all passwords and make a deep copy of the map
	updated := make(map[string]map[string]string)
	for name, section := range s.Sections {
		updated[name] = make(map[string]string)
		for k, v := range section {

			decoded, err := dec.DecodePassword(v)
			if err != nil {
				panic(err)
			}

			encoded, err := enc.EncodePassword(string(decoded))
			if err != nil {
				panic(err)
			}

			updated[name][k] = string(encoded)
		}
	}

	s.Sections = updated

	hashed, err := bcrypt.GenerateFromPassword([]byte(new), bcryptCost)
	if err != nil {
		return err
	}

	s.Passphrase = string(hashed)
	return nil
}

// ListSections retrieves all the sections from the storage.
func (s *Storage) ListSections() []string {
	sections := make([]string, 0)
	for k, _ := range s.Sections {
		sections = append(sections, k)
	}
	return sections
}

// ListPasswords lists all the password contained in the section.
func (s *Storage) ListPasswords(section string) []string {
	res := make([]string, 0)
	sec, ok := s.Sections[section]
	if !ok {
		return res
	}

	for k, _ := range sec {
		res = append(res, k)
	}

	return res
}

// List all lists all the sections and their passwords from the storage.
func (s *Storage) ListAll() map[string][]string {
	res := make(map[string][]string)

	for k, v := range s.Sections {
		res[k] = make([]string, 0)
		for pass, _ := range v {
			res[k] = append(res[k], pass)
		}
	}

	return res
}

// AddSection adds a section to the storage.
func (s *Storage) AddSection(section string) error {
	if _, ok := s.Sections[section]; ok {
		return fmt.Errorf("Section %s already exists", section)
	}

	s.Sections[section] = make(map[string]string)
	return nil
}

// Get retrieves an encrypted password from the storage. If the section or the password do not exist, an error is raised.
func (s *Storage) Get(section string, password string) (string, error) {
	sec, ok := s.Sections[section]

	if !ok {
		return "", fmt.Errorf("Section %s does not exists", section)
	}

	pass, ok := sec[password]
	if !ok {
		return "", fmt.Errorf("Password %s does not exists", section)
	}

	return pass, nil
}

// Set puts a new encrypted password on the storage. Be extra-careful, it does not actually encrypts and encode it.
func (s *Storage) Set(section string, password string, data string) {
	sec, ok := s.Sections[section]
	if !ok {
		sec = make(map[string]string)
		s.Sections[section] = sec
	}

	sec[password] = data

}
