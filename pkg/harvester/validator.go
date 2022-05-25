package harvester

type Validator struct {
	AccountID  string                  `json:"accountId"`
	Status     string                  `json:"status"`
	EraPoints  uint64                  `json:"eraPoints"`
	TotalStake uint64                  `json:"totalStake"`
	OwnStake   uint64                  `json:"ownStake"`
	Commission float64                 `json:"commission"`
	Nominators []string                `json:"nominators"`
	Entity     Entity                  `json:"entity"`
	Account    ValidatorAccountDetails `json:"validatorAccountDetails"`
}

type Entity struct {
	Name      string `json:"name"`
	AccountID string `json:"accountId"`
}

type ValidatorAccountDetails struct {
	Name         string  `json:"name"`
	TotalBalance float64 `json:"totalBalance"`
	Locked       float64 `json:"locked"`
	Bonded       float64 `json:"bonded"`
	Parent       string  `json:"parent"`
	Email        string  `json:"email"`
	Website      string  `json:"website"`
	Twitter      string  `json:"twitter"`
	Riot         string  `json:"riot"`
}
