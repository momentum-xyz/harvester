package substrate

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestValidators(t *testing.T) {
	stashAccounts, _ := mockSh.getStashAccounts()

	t.Run("getCurrentSessionValidators()", func(t *testing.T) {
		validatorAccountIDs, err := mockSh.getCurrentSessionValidators()
		assert.Nil(t, err)
		assert.IsType(t, &validatorAccountIDs[0], &types.AccountID{})
		assert.Greater(t, len(validatorAccountIDs), 0)

	})

	t.Run("processValidators", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(4 * time.Second)
			ctx.Done()
			cancel()
		}()
		err := mockSh.processValidators(ctx, func(err error) {}, mockHarvester.PerformanceMonitorClient, "validators", 1*time.Second)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("updateValidators", func(t *testing.T) {
		err := mockSh.updateValidators(func(err error) {}, mockHarvester.PerformanceMonitorClient, "validators")
		assert.Nil(t, err)
	})

	t.Run("getValidatorPreferences", func(t *testing.T) {
		accountID, _ := StringToAccountId(stashAccounts[0])
		validatorPrefs, err := mockSh.getValidatorPreferences(accountID)
		assert.Nil(t, err)
		assert.IsType(t, validatorPrefs, harvester.ValidatorPreferences{})
	})

	t.Run("getValidatorIdentity", func(t *testing.T) {
		accountID, _ := StringToAccountId(stashAccounts[0])
		identity, err := mockSh.getValidatorIdentity(accountID)
		assert.Nil(t, err)
		assert.IsType(t, identity, harvester.ValidatorInfo{})
	})

	t.Run("getSuperOf", func(t *testing.T) {
		accountID, _ := StringToAccountId(stashAccounts[0])
		superIdentity, err := mockSh.getSuperOf(accountID)
		assert.Nil(t, err)
		assert.IsType(t, superIdentity, harvester.SuperIdentity{})
	})

	t.Run("getSubsOf", func(t *testing.T) {
		accountID, _ := StringToAccountId(stashAccounts[0])
		subs, err := mockSh.getSubsOf(accountID)
		assert.Nil(t, err)
		assert.IsType(t, subs, harvester.SubsIdentity{})
	})

	t.Run("getStashAccounts", func(t *testing.T) {
		accounts, err := mockSh.getStashAccounts()
		assert.Nil(t, err)
		accountID, err := StringToAccountId(accounts[0])
		assert.Nil(t, err)
		assert.IsType(t, accountID, types.AccountID{})
	})

	t.Run("decodeStashKey", func(t *testing.T) {
		key, _ := mockSh.CreateStorageKeyUnsafe("Staking", "Validators")
		stashKeys, _ := mockSh.GetKeysLatest(key)
		res, err := mockSh.decodeStashKey(stashKeys[0])
		assert.Nil(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("getValidatorStatus", func(t *testing.T) {
		accountIds := []types.AccountID{types.NewAccountID([]byte{1}), types.NewAccountID([]byte{2})}
		res := getValidatorStatus(accountIds, types.NewAccountID([]byte{1}))
		assert.Equal(t, res, "active")
		res = getValidatorStatus(accountIds, types.NewAccountID([]byte{3}))
		assert.Equal(t, res, "candidate")
	})

	t.Run("getEraIndex", func(t *testing.T) {
		res, err := mockSh.getEraIndex()
		assert.Nil(t, err)
		assert.Equal(t, reflect.TypeOf(res).String(), reflect.TypeOf(types.NewU32(0)).String())
	})

	t.Run("getEraStakers", func(t *testing.T) {
		accountID, _ := StringToAccountId(stashAccounts[0])
		res, err := mockSh.getEraStakers(accountID)
		assert.Nil(t, err)
		assert.IsType(t, res, types.Exposure{})
	})

	t.Run("getValidatorLockedBalance", func(t *testing.T) {
		accountID, _ := StringToAccountId(stashAccounts[0])
		res, err := mockSh.getValidatorLockedBalance(accountID)
		assert.Nil(t, err)
		assert.Greater(t, len(res), 0)
		assert.IsType(t, res[0], harvester.ValidatorBalancesLocked{})
		assert.Greater(t, res[0].Amount, int64(0))

		res, err = mockSh.getValidatorLockedBalance(types.NewAccountID([]byte{1}))
		assert.Equal(t, len(res), 0)
		assert.Nil(t, err)
	})
}
