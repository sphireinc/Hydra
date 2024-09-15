package main

import (
	"Hydrator/hydra"
	"database/sql"
	"encoding/json"
	"fmt"
)

type Clan struct {
	ClanAcc     string `json:"clan_acc" hydra:"clan_acc"`
	ID          string `json:"id" hydra:"id"`
	PreviousId  string `json:"name" hydra:"previous_id"`
	Description string `json:"age" hydra:"description"`
	Author      string `json:"author" hydra:"author"`
	Comment     string `json:"comment" hydra:"comment"`
	Created     string `json:"created" hydra:"created"`
	Updated     string `json:"updated" hydra:"updated"`
	hydra.Hydratable
}

// Implement the fmt.Stringer interface to output JSON by default
func (c Clan) String() string {
	// Marshal the struct to JSON
	jsonData, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("Error marshaling to JSON: %v", err)
	}
	return string(jsonData)
}

func main() {
	p := &Clan{}
	p.Init(p)

	// test using a publicly available MySQL database
	db, err := sql.Open("mysql", "rfamro:@tcp(mysql-rfam-public.ebi.ac.uk:4497)/Rfam")
	if err != nil {
		panic(err)
	}

	fmt.Println(p.Hydrate(db, map[string]interface{}{"id": "U6"}))
	fmt.Println("res", p)
}
