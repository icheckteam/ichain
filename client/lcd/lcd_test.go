package lcd

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/icheckteam/ichain/x/asset"
	"github.com/icheckteam/ichain/x/identity"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	client "github.com/cosmos/cosmos-sdk/client"
	tests "github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

//_____________________________________________________________________________
// get the account to get the sequence
func getAccount(t *testing.T, port string, addr sdk.AccAddress) auth.Account {
	res, body := Request(t, port, "GET", fmt.Sprintf("/accounts/%s", addr), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var acc auth.Account
	err := cdc.UnmarshalJSON([]byte(body), &acc)
	require.Nil(t, err)
	return acc
}

func TestCreateAsset(t *testing.T) {

	name, password := "test", "1234567890"
	addr, seed := CreateAddr(t, "test", password, GetKB(t))
	cleanup, _, port := InitializeTestLCD(t, 1, []sdk.AccAddress{addr})
	defer cleanup()

	// CreateAsset tests
	resultTx := doCreateAsset(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)

	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	asset := getAsset(t, port, "test")
	assert.Equal(t, asset.ID, "test")

	// UpdateProperties tests
	resultTx = doUpdateProperties(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	// AddMaterials Tests
	resultTx = doAddMaterials(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	// AddQuantity
	resultTx = doAddQuantity(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	// SubtractQuantity
	resultTx = doSubtractQuantity(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	// doFinalizeAsset
	resultTx = doFinalizeAsset(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)
}

func TestCreateProposal(t *testing.T) {
	name, password := "test", "1234567890"
	name2, password2 := "test2", "1234567890"
	addr, seed := CreateAddr(t, name, password, GetKB(t))
	addr2, seed2 := CreateAddr(t, name2, password2, GetKB(t))
	cleanup, _, port := InitializeTestLCD(t, 1, []sdk.AccAddress{addr, addr2})
	defer cleanup()

	// CreateAsset tests
	resultTx := doCreateAsset(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)

	// CreateProposal tests
	resultTx = doCreateProposal(t, port, seed, name, password, addr, addr2)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	// getProposals
	proposals := getProposals(t, port)
	assert.Equal(t, len(proposals), 1)

	// AnswerProposal tests
	resultTx = doAnswerProposal(t, port, seed2, name2, password2, addr2, addr2)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	asset := getAsset(t, port, "test")
	assert.Equal(t, len(asset.Reporters), 1)

	// RevokeReporter tests
	resultTx = doRevokeReporter(t, port, seed2, name, password, addr, addr2)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	asset = getAsset(t, port, "test")
	assert.Equal(t, len(asset.Reporters), 0)
}

func doCreateAsset(t *testing.T, port, seed, name, password string, addr sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()

	chainID := viper.GetString(client.FlagChainID)

	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": %d,
			"sequence": %d,
			"gas": 10000,
			"chain_id": "%s"
		},
		"name": "test",
		"asset_id": "test",
		"quantity": "100",
		"unit": "kg",
		"properties": [
			{"name": "size", "type": 4, "number_value": 50}
		]
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", "/assets", jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx

}

func doUpdateProperties(t *testing.T, port, seed, name, password string, addr sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)
	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": %d,
			"sequence": %d,
			"gas": 10000,
			"chain_id": "%s"
		},
		"properties": [
			{"name": "size", "type": 4, "number_value": 50}
		]
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", "/assets/test/properties", jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx

}

func doAddMaterials(t *testing.T, port, seed, name, password string, addr sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)
	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": %d,
			"sequence": %d,
			"gas": 10000,
			"chain_id": "%s"
		},
		"amount": [
			{"denom": "test", "amount": "5"}
		]
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", "/assets/test/materials", jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func doAddQuantity(t *testing.T, port, seed, name, password string, addr sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)
	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": %d,
			"sequence": %d,
			"gas": 10000,
			"chain_id": "%s"
		},
		"quantity": "5"
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", "/assets/test/add", jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func doSubtractQuantity(t *testing.T, port, seed, name, password string, addr sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)
	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": %d,
			"sequence": %d,
			"gas": 10000,
			"chain_id": "%s"
		},
		"quantity": "5"
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", "/assets/test/subtract", jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func doFinalizeAsset(t *testing.T, port, seed, name, password string, addr sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)
	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": %d,
			"sequence": %d,
			"gas": 10000,
			"chain_id": "%s"
		}
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", "/assets/test/finalize", jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func doCreateProposal(t *testing.T, port, seed, name, password string, addr, recipient sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)

	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": %d,
			"sequence": %d,
			"gas": 10000,
			"chain_id": "%s"
		},
		"recipient": "%s",
		"properties": ["size"],
		"role": 1
	}`, name, password, accnum, sequence, chainID, recipient))

	res, body := Request(t, port, "POST", "/assets/test/proposals", jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func doAnswerProposal(t *testing.T, port, seed, name, password string, addr, recipient sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)

	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": %d,
			"sequence": %d,
			"gas": 10000,
			"chain_id": "%s"
		},
		"response": 1
	}`, name, password, accnum, sequence, chainID))
	res, body := Request(t, port, "POST", fmt.Sprintf("/assets/test/proposals/%s/answer", recipient), jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func doRevokeReporter(t *testing.T, port, seed, name, password string, addr, recipient sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)
	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": %d,
			"sequence": %d,
			"gas": 10000,
			"chain_id": "%s"
		}
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", fmt.Sprintf("/assets/test/reporters/%s/revoke", recipient), jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func getAsset(t *testing.T, port string, assetID string) asset.Asset {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", fmt.Sprintf("/assets/%s", assetID), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var a asset.Asset
	err := cdc.UnmarshalJSON([]byte(body), &a)
	require.Nil(t, err)
	return a
}

func getProposals(t *testing.T, port string) []asset.Proposal {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", fmt.Sprintf("/assets/%s/proposals", "test"), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var proposals []asset.Proposal
	err := cdc.UnmarshalJSON([]byte(body), &proposals)
	require.Nil(t, err)
	return proposals
}

// Test Identity Module
// -------------------------------------------------------------------------------------------------

func TestAddTrust(t *testing.T) {
	name, password := "test", "1234567890"
	addr, seed := CreateAddr(t, name, password, GetKB(t))
	name2, _ := "test2", "1234567890"
	addr2, _ := CreateAddr(t, name2, password, GetKB(t))
	cleanup, _, port := InitializeTestLCD(t, 1, []sdk.AccAddress{addr, addr2})
	defer cleanup()

	// AddTrust tests
	resultTx := doAddTrust(t, port, seed, name, password, addr, addr2)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	trusts := getTrusts(t, port, addr)
	assert.Equal(t, len(trusts), 1)

}

func TestCreateIdentity(t *testing.T) {
	name, password := "test", "1234567890"
	addr, seed := CreateAddr(t, name, password, GetKB(t))
	cleanup, _, port := InitializeTestLCD(t, 1, []sdk.AccAddress{addr})
	defer cleanup()

	// AddTrust tests
	resultTx := doCreateIdentity(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	identities := getIdentsByAccount(t, port, addr)
	assert.Equal(t, identities[0].ID, int64(1))

}

func TestAddCerts(t *testing.T) {
	name, password := "test", "1234567890"
	addr, seed := CreateAddr(t, name, password, GetKB(t))
	cleanup, _, port := InitializeTestLCD(t, 1, []sdk.AccAddress{addr})
	defer cleanup()

	// CreateIdentity
	resultTx := doCreateIdentity(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)

	// AddCerts tests
	resultTx = doAddCerts(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	certs := getCertsByIdentity(t, port)
	assert.Equal(t, len(certs), 2)

	ident := getClaimedIdentity(t, port, addr)
	assert.Equal(t, ident.ID, int64(1))
}

func doAddTrust(t *testing.T, port, seed, name, password string, addr, trusting sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)
	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": %d,
			"sequence": %d,
			"gas": 10000,
			"chain_id": "%s"
		},
		"trust": true
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", fmt.Sprintf("/accounts/%s/trusts", trusting), jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func doCreateIdentity(t *testing.T, port, seed, name, password string, addr sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)
	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": %d,
			"sequence": %d,
			"gas": 10000,
			"chain_id": "%s"
		}
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", "/identities", jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func doAddCerts(t *testing.T, port, seed, name, password string, addr sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)
	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": %d,
			"sequence": %d,
			"gas": 10000,
			"chain_id": "%s"
		},
		"values": [
			{
				"property": "company",
				"type": "demo",
				"data": {
					"demo": "1212"
				},
				"confidence": true
			},
			{
				"property": "owner",
				"confidence": true
			}
		]
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", fmt.Sprintf("/identities/1/certs"), jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func getTrusts(t *testing.T, port string, trustor sdk.AccAddress) []identity.Trust {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", fmt.Sprintf("/accounts/%s/trusts", trustor), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var trusts []identity.Trust
	err := cdc.UnmarshalJSON([]byte(body), &trusts)
	require.Nil(t, err)
	return trusts
}

func getIdentsByAccount(t *testing.T, port string, addr sdk.AccAddress) []identity.Identity {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", fmt.Sprintf("/accounts/%s/identities", addr), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var identities []identity.Identity
	err := cdc.UnmarshalJSON([]byte(body), &identities)
	require.Nil(t, err)
	return identities
}

func getCertsByIdentity(t *testing.T, port string) []identity.Cert {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", "/identities/1/certs", nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var certs []identity.Cert
	err := cdc.UnmarshalJSON([]byte(body), &certs)
	require.Nil(t, err)
	return certs
}

func getClaimedIdentity(t *testing.T, port string, addr sdk.AccAddress) identity.Identity {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", fmt.Sprintf("/accounts/%s/claimed", addr), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var ident identity.Identity
	err := cdc.UnmarshalJSON([]byte(body), &ident)
	require.Nil(t, err)
	return ident
}
