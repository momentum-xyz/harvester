package main

import (
	"encoding/json"
	"fmt"
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
)

func main() {
	validator := harvester.Validator{
		AccountID:  "0x11111111111111111111111",
		Status:     "active",
		EraPoints:  100,
		TotalStake: 10,
		OwnStake:   1,
		Commission: 3.0,
		Nominators: []string{"0x2222222222222222222222", "0x3333333333333333"},
		Entity: harvester.Entity{
			Name:      "SUPER VALIDATORS",
			AccountID: "0x44444444444444444444444444",
		},
		Account: harvester.ValidatorAccountDetails{
			Name:         "SUPER VALIDATOR #1",
			TotalBalance: 100,
			Locked:       10,
			Bonded:       13,
			Parent:       "0x44444444444444444444444444",
			Email:        "super@supervalidators.com",
			Website:      "supervalidators.com",
			Twitter:      "@supervalidators",
			Riot:         "@supervalidators",
		},
	}
	j, _ := json.MarshalIndent(validator, "", "    ")
	fmt.Println(string(j))
}
