package tests

import "Hydrator/hydra"

type Person struct {
	Id          int    `json:"id" hydra:"id"`
	FirstName   string `json:"first_name" hydra:"first_name"`
	LastName    string `json:"last_name" hydra:"last_name"`
	Sex         string `json:"sex" hydra:"sex"`
	DateOfBirth string `json:"date_of_birth" hydra:"date_of_birth"`
	Created     string `json:"created" hydra:"created"`
	Updated     string `json:"updated" hydra:"updated"`
	Deleted     string `json:"deleted" hydra:"deleted"`
	hydra.Hydratable
}

type Address struct {
	Id         int    `json:"id" hydra:"id"`
	UserId     int    `json:"user_id" hydra:"user_id"`
	Address1   string `json:"address_1" hydra:"address_1"`
	Address2   string `json:"address_2" hydra:"address_2"`
	City       string `json:"city" hydra:"city"`
	State      string `json:"state" hydra:"state"`
	Province   string `json:"province" hydra:"province"`
	PostalCode string `json:"postal_code" hydra:"postal_code"`
	Country    string `json:"country" hydra:"country"`
	Created    string `json:"created" hydra:"created"`
	Updated    string `json:"updated" hydra:"updated"`
	Deleted    string `json:"deleted" hydra:"deleted"`
	hydra.Hydratable
}
