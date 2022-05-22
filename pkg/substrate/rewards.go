package substrate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
	"github.com/OdysseyMomentumExperience/harvester/pkg/log"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/vedhavyas/go-subkey"
)

type Reward struct {
	Era         uint32            `json:"era"`
	TotalPoints uint32            `json:"totalPoints"`
	Rewards     map[string]uint32 `json:"rewards"`
}

type EraRewardPoints struct {
	Total      types.U32
	Individual []struct {
		AccountID    types.AccountID
		RewardPoints types.U32
	}
}

func (sh *SubstrateHarvester) ProcessErasRewardPoints(ctx context.Context, fn harvester.ErrorHandler) error {

	ticker := time.NewTicker(60 * time.Second)

	defer ticker.Stop()

	for {
		var err error
		select {
		case <-ticker.C:
			sh.publishErasRewardPoints()
		case <-ctx.Done():
			return ctx.Err()
		}
		if err != nil {
			fn(err)
		}
	}

}

func (sh *SubstrateHarvester) publishErasRewardPoints() error {
	activeEra, err := sh.GetActiveEra()
	if err != nil {
		return errors.Wrap(err, "error while fetching active era")
	}

	activeEraDepth, err := sh.GetEraDepth(activeEra)

	if err != nil {
		return errors.Wrap(err, "error while fetching active era depth")
	}

	rewards, err := sh.GetEraReward(activeEraDepth)
	if err != nil {
		return err
	}

	rewardsMap := make(map[string]uint32)
	for _, v := range rewards.Individual {
		address, err := subkey.SS58Address(v.AccountID[:], 2)
		if err != nil {
			return err
		}

		rewardsMap[address] = uint32(v.RewardPoints)
	}

	rewardsJson, err := json.Marshal(Reward{
		Era:         activeEra,
		TotalPoints: uint32(rewards.Total),
		Rewards:     rewardsMap,
	})
	if err != nil {
		return err
	}

	log.Logln(0, fmt.Sprintf("%s - Publishing reward event for era %d", sh.cfg.Name, activeEra))
	err = sh.publisher.Publish(fmt.Sprintf("harvester/%s/reward-event", sh.cfg.Name), string(rewardsJson))
	if err != nil {
		return err
	}
	return nil
}

func (sh *SubstrateHarvester) GetEraReward(era []byte) (EraRewardPoints, error) {

	noRewards := EraRewardPoints{}

	key, err := sh.GetStorageDataKey("Staking", "ErasRewardPoints", era)
	if err != nil {
		return noRewards, errors.Wrapf(err, "error while fetching storage key for erasRewardPoints for era %v", spew.Sdump(era))
	}

	var rewards EraRewardPoints
	err = sh.GetStorageLatest(key, &rewards)
	if err != nil {
		return noRewards, errors.Wrapf(err, "error while fetching storage erasRewardPoints for key %v", spew.Sdump(key))
	}

	return rewards, nil
}

// ## TODO FOR REWARD & SLASH EVENTS (upcoming story) AND GOOD REFERENCE OF EXTENDING TO CUSTOM EVENTS
// func (sh *SubstrateHarvester) GetRewardEvents() {
// 	// Subscribe to system events via storage
// 	key, err := types.CreateStorageKey(sh.metadata, "System", "Events", nil, nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	sub, err := sh.api.RPC.State.SubscribeStorageRaw([]types.StorageKey{key})
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer sub.Unsubscribe()

// 	// outer for loop for subscription notifications
// 	for {
// 		set := <-sub.Chan()
// 		// inner loop for the changes within one of those notifications
// 		for _, chng := range set.Changes {
// 			if !types.Eq(chng.StorageKey, key) || !chng.HasStorageData {
// 				// skip, we are only interested in events with content
// 				continue
// 			}

// 			// Decode the event records
// 			events := EventRecords{}
// 			err = EventRecordsRaw(chng.StorageData).DecodeEventRecords(sh.metadata, &events)
// 			if err != nil {
// 				//fmt.Println("## ERROR ##")
// 				//fmt.Println(err)
// 				continue
// 			}

// 			// Show what we are busy with
// 			for _, e := range events.Staking_Reward {
// 				fmt.Printf("\tStaking:Reward:: (phase=%#v)\n", e.Phase)
// 				fmt.Printf("\t\t%v\n", e.Amount)
// 				fmt.Println(e)
// 			}
// 			for _, e := range events.Staking_Slash {
// 				fmt.Printf("\tStaking:Slash:: (phase=%#v)\n", e.Phase)
// 				fmt.Printf("\t\t%#x%v\n", e.AccountID, e.Balance)
// 				fmt.Println(e)
// 			}

// 		}
// 	}
// }

type EventRecordsRaw []byte

// Override of DecodeEventRecords from gsrpc pkg
func (e EventRecordsRaw) DecodeEventRecords(m *types.Metadata, t interface{}) error {
	log.Logln(2, fmt.Sprintf("will decode event records from raw hex: %#x", e))

	// ensure t is a pointer
	ttyp := reflect.TypeOf(t)
	if ttyp.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer, but is " + fmt.Sprint(ttyp))
	}
	// ensure t is not a nil pointer
	tval := reflect.ValueOf(t)
	if tval.IsNil() {
		return errors.New("target is a nil pointer")
	}
	val := tval.Elem()
	typ := val.Type()
	// ensure val can be set
	if !val.CanSet() {
		return fmt.Errorf("unsettable value %v", typ)
	}
	// ensure val points to a struct
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("target must point to a struct, but is " + fmt.Sprint(typ))
	}

	decoder := scale.NewDecoder(bytes.NewReader(e))

	// determine number of events
	n, err := decoder.DecodeUintCompact()
	if err != nil {
		return err
	}

	log.Logln(2, fmt.Sprintf("found %v events", n))

	// iterate over events
	for i := uint64(0); i < n.Uint64(); i++ {
		log.Logln(2, fmt.Sprintf("decoding event #%v", i))

		// decode Phase
		phase := types.Phase{}
		err := decoder.Decode(&phase)
		if err != nil {
			return fmt.Errorf("unable to decode Phase for event #%v: %v", i, err)
		}

		// decode EventID
		id := types.EventID{}
		err = decoder.Decode(&id)
		if err != nil {
			return fmt.Errorf("unable to decode EventID for event #%v: %v", i, err)
		}

		log.Logln(1, fmt.Sprintf("event #%v has EventID %v", i, id))

		// ask metadata for method & event name for event
		moduleName, eventName, err := m.FindEventNamesForEventID(id)
		// moduleName, eventName, err := "System", "ExtrinsicSuccess", nil
		if err != nil {
			fmt.Printf("unable to find event with EventID %v in metadata for event #%v: %s\n", id, i, err)
			break
		}

		log.Logln(0, fmt.Sprintf("event #%v is in module %v with event name %v", i, moduleName, eventName))

		// check whether name for eventID exists in t
		field := val.FieldByName(fmt.Sprintf("%v_%v", moduleName, eventName))
		if !field.IsValid() {
			return fmt.Errorf("unable to find field %v_%v for event #%v with EventID %v", moduleName, eventName, i, id)
		}

		// create a pointer to with the correct type that will hold the decoded event
		holder := reflect.New(field.Type().Elem())

		// ensure first field is for Phase, last field is for Topics
		numFields := holder.Elem().NumField()
		if numFields < 2 {
			return fmt.Errorf("expected event #%v with EventID %v, field %v_%v to have at least 2 fields "+
				"(for Phase and Topics), but has %v fields", i, id, moduleName, eventName, numFields)
		}
		phaseField := holder.Elem().FieldByIndex([]int{0})
		if phaseField.Type() != reflect.TypeOf(phase) {
			return fmt.Errorf("expected the first field of event #%v with EventID %v, field %v_%v to be of type "+
				"types.Phase, but got %v", i, id, moduleName, eventName, phaseField.Type())
		}
		topicsField := holder.Elem().FieldByIndex([]int{numFields - 1})
		if topicsField.Type() != reflect.TypeOf([]types.Hash{}) {
			return fmt.Errorf("expected the last field of event #%v with EventID %v, field %v_%v to be of type "+
				"[]types.Hash for Topics, but got %v", i, id, moduleName, eventName, topicsField.Type())
		}

		// set the phase we decoded earlier
		phaseField.Set(reflect.ValueOf(phase))

		// set the remaining fields
		for j := 1; j < numFields; j++ {
			err = decoder.Decode(holder.Elem().FieldByIndex([]int{j}).Addr().Interface())
			if err != nil {
				return fmt.Errorf("unable to decode field %v event #%v with EventID %v, field %v_%v: %v", j, i, id, moduleName,
					eventName, err)
			}
		}

		// add the decoded event to the slice
		field.Set(reflect.Append(field, holder.Elem()))

		log.Logln(2, fmt.Sprintf("decoded event #%v", i))
	}
	return nil
}
