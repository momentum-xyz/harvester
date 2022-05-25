package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

// custom event types
type EventParachainsCandidateBacked struct {
	Phase                                   types.Phase
	PolkadotPrimitivesV1CandidateDescriptor PolkadotPrimitivesV1CandidateDescriptor
	HeadData1                               types.Bytes
	HeadData2                               types.U32
	HeadData3                               types.U32
	Topics                                  []types.Hash
}

type EventParaInclusionCandidateIncluded struct {
	Phase                                   types.Phase
	PolkadotPrimitivesV1CandidateDescriptor PolkadotPrimitivesV1CandidateDescriptor
	HeadData1                               types.Bytes
	HeadData2                               types.U32
	HeadData3                               types.U32
	Topics                                  []types.Hash
}

type PolkadotPrimitivesV1CandidateDescriptor struct {
	ParaID                      types.U32
	RelayParent                 types.H256
	Collator                    PolkadotPrimitivesV0CollatorAppPublic
	PersistedValidationDataHash types.H256
	PovHash                     types.H256
	ErasureRoot                 types.H256
	Signature                   PolkadotPrimitivesV0CollatorAppSignature
	Parahead                    types.H256
	ValidationCodeHash          types.H256
}

type PolkadotPrimitivesV0CollatorAppPublic struct {
	types.U8
}

type PolkadotPrimitivesV0CollatorAppSignature struct {
	types.U8
}

// Override of event records from types in order to complement events
type EventRecords struct {
	types.EventRecords
	ParaInclusion_CandidateIncluded []EventParaInclusionCandidateIncluded
	ParaInclusion_CandidateBacked   []EventParachainsCandidateBacked
}
