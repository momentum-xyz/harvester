package harvester

import (
	"encoding/json"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type Validator struct {
	AccountID  string                    `json:"accountId"`
	Name       string                    `json:"name"`
	Commission float64                   `json:"commission"`
	Status     string                    `json:"status"`
	Balance    string                    `json:"balance"`
	Reserved   string                    `json:"reserved"`
	Locked     []ValidatorBalancesLocked `json:"locked"`
	OwnStake   string                    `json:"ownStake"`
	TotalStake string                    `json:"totalStake"`
	Identity   ValidatorInfo             `json:"identity"`
	Nominators []Nominator               `json:"nominators"`
	Parent     Parent                    `json:"parent"`
	Children   []string                  `json:"children"`
	Hash       string                    `json:"hash"`
	Chain      string                    `json:"chain"`
}

type ValidatorBalancesLocked struct {
	ID      string
	Amount  int64
	Reasons string
}

type ValidatorInfo struct {
	Display        string               `json:"display"`
	Legal          string               `json:"legal"`
	Web            string               `json:"web"`
	Riot           string               `json:"riot"`
	Email          string               `json:"email"`
	PgpFingerprint string               `json:"pgpFingerprint"`
	Image          string               `json:"image"`
	Twitter        string               `json:"twitter"`
	Parent         string               `json:"parent"`
	DisplayParent  string               `json:"displayParent"`
	Judgements     []ValidatorJudgement `json:"judgements"`
}

type ValidatorJudgement struct {
	Index types.U32               `json:"index"`
	Value PalletIdentityJudgement `json:"value"`
}

func (j *ValidatorJudgement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&[2]string{
		fmt.Sprint(j.Index),
		j.Value.String(),
	})
}

type PalletIdentityJudgement uint8

type Nominator struct {
	Address string `json:"address"`
	Stake   string `json:"stake"`
}

type Parent struct {
	Name      string `json:"name"`
	AccountID string `json:"accountId"`
}

type ValidatorPreferences struct {
	Commission types.UCompact
	Blocked    types.Bool
}

type PalletBalancesReasons uint8

const (
	Fee PalletBalancesReasons = iota
	Misc
	All
)

func (a PalletBalancesReasons) String() string {
	reasons := [...]string{"fee", "misc", "all"}
	if int(a) < len(reasons) {
		return reasons[a]
	}

	return "fee"
}

type BalancesLocked struct {
	ID      types.Bytes8
	Amount  types.U128
	Reasons PalletBalancesReasons
}

const (
	IsUnknown PalletIdentityJudgement = iota
	IsFeePaid
	IsReasonable
	IsKnownGood
	IsOutOfDate
	IsLowQuality
	IsErroneous
)

func (s PalletIdentityJudgement) String() string {
	judgements := [...]string{
		"Unknown",
		"Fee Paid",
		"Reasonable",
		"Known Good",
		"Out Of Date",
		"Low Quality",
		"Erroneous",
	}

	if int(s) < len(judgements) {
		return judgements[s]
	}

	return "Unknown"
}

type PalletIdentityInfo struct {
	Additional []struct {
		Key   types.Data
		Value types.Data
	}
	Rest types.Data
}

type AccountInfo struct {
	Nonce       types.U32
	Consumers   types.U32
	Providers   types.U32
	Sufficients types.U32
	Data        struct {
		Free       types.U128
		Reserved   types.U128
		MiscFrozen types.U128
		FeeFrozen  types.U128
	}
}

type Identity struct {
	Judgements []ValidatorJudgement
	Deposit    types.U128
	Info       PalletIdentityInfo
}

type SuperIdentity struct {
	Account types.AccountID
	Name    types.Data
}

type SubsIdentity struct {
	Deposit  types.U128
	Accounts []types.AccountID
}

type IdentityByte types.Data

func (d *IdentityByte) Decode() string {
	var result string
	if len(*d) < 1 {
		return result
	}

	first := (*d)[0]

	switch {
	case first == 0:
		*d = (*d)[1:]
		result = "None"
	case first >= 1 && first <= 33:
		store := []byte{}

		for i := byte(1); i < first; i++ {
			store = append(store, (*d)[i])
		}

		*d = (*d)[first:]
		result = string(store)
	}

	return result
}
