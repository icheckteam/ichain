package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

var _ sdk.Account = (*AppAccount)(nil)

// AppAccount Custom extensions for this application.  This is just an example of
// extending auth.BaseAccount with custom fields.
//
// This is compatible with the stock auth.AccountStore, since
// auth.AccountStore uses the flexible go-wire library.
type AppAccount struct {
	auth.BaseAccount
	Name       string `json:"name"`
	Identities Identities
	Assets     sdk.Coins `json:"assets"`
}

// nolint
func (acc AppAccount) GetName() string                      { return acc.Name }
func (acc *AppAccount) SetName(name string)                 { acc.Name = name }
func (acc AppAccount) GetIdentities() Identities            { return acc.Identities }
func (acc *AppAccount) SetIdentities(identities Identities) { acc.Identities = identities }
func (acc AppAccount) GetAssets() sdk.Coins                 { return acc.Assets }
func (acc *AppAccount) SetAssets(assets sdk.Coins)          { acc.Assets = assets }

// GetAccountDecoder Get the AccountDecoder function for the custom AppAccount
func GetAccountDecoder(cdc *wire.Codec) sdk.AccountDecoder {
	return func(accBytes []byte) (res sdk.Account, err error) {
		if len(accBytes) == 0 {
			return nil, sdk.ErrTxDecode("accBytes are empty")
		}
		acct := new(AppAccount)
		err = cdc.UnmarshalBinary(accBytes, &acct)
		if err != nil {
			panic(err)
		}
		return acct, err
	}
}

//___________________________________________________________________________________

// GenesisState State to Unmarshal
type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Name    string      `json:"name"`
	Address sdk.Address `json:"address"`
	Coins   sdk.Coins   `json:"coins"`
}

// NewGenesisAccount new genesis account
func NewGenesisAccount(aa *AppAccount) *GenesisAccount {
	return &GenesisAccount{
		Name:    aa.Name,
		Address: aa.Address,
		Coins:   aa.Coins,
	}
}

// ToAppAccount convert GenesisAccount to AppAccount
func (ga *GenesisAccount) ToAppAccount() (acc *AppAccount, err error) {
	baseAcc := auth.BaseAccount{
		Address: ga.Address,
		Coins:   ga.Coins,
	}
	return &AppAccount{
		BaseAccount: baseAcc,
		Name:        ga.Name,
	}, nil
}
