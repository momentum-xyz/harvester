// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/OdysseyMomentumExperience/harvester/ent/validator"
	"github.com/OdysseyMomentumExperience/harvester/pkg/harvester"
)

// Validator is the model entity for the Validator schema.
type Validator struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// AccountID holds the value of the "account_id" field.
	AccountID string `json:"account_id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Commission holds the value of the "commission" field.
	Commission float64 `json:"commission,omitempty"`
	// Status holds the value of the "status" field.
	Status string `json:"status,omitempty"`
	// Balance holds the value of the "balance" field.
	Balance string `json:"balance,omitempty"`
	// Reserved holds the value of the "reserved" field.
	Reserved string `json:"reserved,omitempty"`
	// Locked holds the value of the "locked" field.
	Locked []harvester.ValidatorBalancesLocked `json:"locked,omitempty"`
	// OwnStake holds the value of the "own_stake" field.
	OwnStake string `json:"own_stake,omitempty"`
	// TotalStake holds the value of the "total_stake" field.
	TotalStake string `json:"total_stake,omitempty"`
	// Identity holds the value of the "identity" field.
	Identity harvester.ValidatorInfo `json:"identity,omitempty"`
	// Nominators holds the value of the "nominators" field.
	Nominators []harvester.Nominator `json:"nominators,omitempty"`
	// Parent holds the value of the "parent" field.
	Parent harvester.Parent `json:"parent,omitempty"`
	// Children holds the value of the "children" field.
	Children []string `json:"children,omitempty"`
	// Hash holds the value of the "hash" field.
	Hash string `json:"hash,omitempty"`
	// Chain holds the value of the "chain" field.
	Chain string `json:"chain,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Validator) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case validator.FieldLocked, validator.FieldIdentity, validator.FieldNominators, validator.FieldParent, validator.FieldChildren:
			values[i] = new([]byte)
		case validator.FieldCommission:
			values[i] = new(sql.NullFloat64)
		case validator.FieldID:
			values[i] = new(sql.NullInt64)
		case validator.FieldAccountID, validator.FieldName, validator.FieldStatus, validator.FieldBalance, validator.FieldReserved, validator.FieldOwnStake, validator.FieldTotalStake, validator.FieldHash, validator.FieldChain:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Validator", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Validator fields.
func (v *Validator) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case validator.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			v.ID = int(value.Int64)
		case validator.FieldAccountID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field account_id", values[i])
			} else if value.Valid {
				v.AccountID = value.String
			}
		case validator.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				v.Name = value.String
			}
		case validator.FieldCommission:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field commission", values[i])
			} else if value.Valid {
				v.Commission = value.Float64
			}
		case validator.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				v.Status = value.String
			}
		case validator.FieldBalance:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field balance", values[i])
			} else if value.Valid {
				v.Balance = value.String
			}
		case validator.FieldReserved:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field reserved", values[i])
			} else if value.Valid {
				v.Reserved = value.String
			}
		case validator.FieldLocked:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field locked", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &v.Locked); err != nil {
					return fmt.Errorf("unmarshal field locked: %w", err)
				}
			}
		case validator.FieldOwnStake:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field own_stake", values[i])
			} else if value.Valid {
				v.OwnStake = value.String
			}
		case validator.FieldTotalStake:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field total_stake", values[i])
			} else if value.Valid {
				v.TotalStake = value.String
			}
		case validator.FieldIdentity:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field identity", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &v.Identity); err != nil {
					return fmt.Errorf("unmarshal field identity: %w", err)
				}
			}
		case validator.FieldNominators:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field nominators", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &v.Nominators); err != nil {
					return fmt.Errorf("unmarshal field nominators: %w", err)
				}
			}
		case validator.FieldParent:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field parent", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &v.Parent); err != nil {
					return fmt.Errorf("unmarshal field parent: %w", err)
				}
			}
		case validator.FieldChildren:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field children", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &v.Children); err != nil {
					return fmt.Errorf("unmarshal field children: %w", err)
				}
			}
		case validator.FieldHash:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field hash", values[i])
			} else if value.Valid {
				v.Hash = value.String
			}
		case validator.FieldChain:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field chain", values[i])
			} else if value.Valid {
				v.Chain = value.String
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Validator.
// Note that you need to call Validator.Unwrap() before calling this method if this Validator
// was returned from a transaction, and the transaction was committed or rolled back.
func (v *Validator) Update() *ValidatorUpdateOne {
	return (&ValidatorClient{config: v.config}).UpdateOne(v)
}

// Unwrap unwraps the Validator entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (v *Validator) Unwrap() *Validator {
	tx, ok := v.config.driver.(*txDriver)
	if !ok {
		panic("ent: Validator is not a transactional entity")
	}
	v.config.driver = tx.drv
	return v
}

// String implements the fmt.Stringer.
func (v *Validator) String() string {
	var builder strings.Builder
	builder.WriteString("Validator(")
	builder.WriteString(fmt.Sprintf("id=%v", v.ID))
	builder.WriteString(", account_id=")
	builder.WriteString(v.AccountID)
	builder.WriteString(", name=")
	builder.WriteString(v.Name)
	builder.WriteString(", commission=")
	builder.WriteString(fmt.Sprintf("%v", v.Commission))
	builder.WriteString(", status=")
	builder.WriteString(v.Status)
	builder.WriteString(", balance=")
	builder.WriteString(v.Balance)
	builder.WriteString(", reserved=")
	builder.WriteString(v.Reserved)
	builder.WriteString(", locked=")
	builder.WriteString(fmt.Sprintf("%v", v.Locked))
	builder.WriteString(", own_stake=")
	builder.WriteString(v.OwnStake)
	builder.WriteString(", total_stake=")
	builder.WriteString(v.TotalStake)
	builder.WriteString(", identity=")
	builder.WriteString(fmt.Sprintf("%v", v.Identity))
	builder.WriteString(", nominators=")
	builder.WriteString(fmt.Sprintf("%v", v.Nominators))
	builder.WriteString(", parent=")
	builder.WriteString(fmt.Sprintf("%v", v.Parent))
	builder.WriteString(", children=")
	builder.WriteString(fmt.Sprintf("%v", v.Children))
	builder.WriteString(", hash=")
	builder.WriteString(v.Hash)
	builder.WriteString(", chain=")
	builder.WriteString(v.Chain)
	builder.WriteByte(')')
	return builder.String()
}

// Validators is a parsable slice of Validator.
type Validators []*Validator

func (v Validators) config(cfg config) {
	for _i := range v {
		v[_i].config = cfg
	}
}
