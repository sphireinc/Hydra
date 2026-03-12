//go:build integration
// +build integration

package functional_tests

import "github.com/sphireinc/Hydra/hydra"

type Person struct {
	ID        int    `json:"id" hydra:"id"`
	FirstName string `json:"first_name" hydra:"first_name"`
	LastName  string `json:"last_name" hydra:"last_name"`
	Sex       string `json:"sex" hydra:"sex"`
	hydra.Hydratable
}

func (Person) HydraTableName() string {
	return "Person"
}

type Address struct {
	ID         int    `json:"id" hydra:"id"`
	UserID     int    `json:"user_id" hydra:"user_id"`
	Address1   string `json:"address_1" hydra:"address_1"`
	Address2   string `json:"address_2" hydra:"address_2"`
	City       string `json:"city" hydra:"city"`
	State      string `json:"state" hydra:"state"`
	PostalCode string `json:"postal_code" hydra:"postal_code"`
	Country    string `json:"country" hydra:"country"`
	hydra.Hydratable
}

func (Address) HydraTableName() string {
	return "Addresses"
}
