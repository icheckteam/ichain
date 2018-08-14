package lcd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/icheckteam/ichain/x/asset"
	"github.com/icheckteam/ichain/x/identity"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func TestCreateAsset(t *testing.T) {

	name, password := "test", "1234567890"
	addr, seed := CreateAddr(t, "test", password, GetKeyBase(t))
	cleanup, _, port := InitializeTestLCD(t, 1, []sdk.AccAddress{addr})
	defer cleanup()

	// CreateAsset tests
	resultTx := doCreateAsset(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)

	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	asset := getAsset(t, port, "test")
	assert.Equal(t, asset.ID, "test")
	assert.Equal(t, len(asset.Properties), 1)

	records := getAssetsByAccount(t, port, asset.Owner)
	assert.Equal(t, len(records), 1)

	properties := getTxsProperties(t, port)
	assert.Equal(t, len(properties), 1)

	owners := getRecordOwners(t, port)
	assert.Equal(t, len(owners), 1)

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

	txs := getTxsTransferMaterials(t, port)
	assert.Equal(t, len(txs), 1)

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
	addr, seed := CreateAddr(t, name, password, GetKeyBase(t))
	addr2, seed2 := CreateAddr(t, name2, password2, GetKeyBase(t))
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

	proposalsOwner := getProposalByOwner(t, port, addr2)
	assert.Equal(t, len(proposalsOwner), 1)

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
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
			"chain_id": "%s"
		},
		"name": "test",
		"asset_id": "test",
		"quantity": "100",
		"unit": "kg",
		"properties": [
			{"name": "size", "type": "4", "number_value": "50"}
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
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
			"chain_id": "%s"
		},
		"properties": [
			{"name": "size", "type": "4", "number_value": "50"}
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

func getTxsTransferMaterials(t *testing.T, port string) []asset.HistoryTransferMaterial {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", fmt.Sprintf("/assets/%s/materials/history", "test"), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var history []asset.HistoryTransferMaterial
	err := json.Unmarshal([]byte(body), &history)
	require.Nil(t, err)
	return history
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
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
			"chain_id": "%s"
		},
		"amount": [
			{"record_id": "test", "amount": "5"}
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
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
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
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
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
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
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
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
			"chain_id": "%s"
		},
		"recipient": "%s",
		"properties": ["size"],
		"role": "1"
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
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
			"chain_id": "%s"
		},
		"response": "1",
		"role": "1"
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
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
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

func getAssetsByAccount(t *testing.T, port string, owner sdk.AccAddress) []*asset.RecordOutput {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", fmt.Sprintf("/accounts/%s/assets", owner.String()), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	record := []*asset.RecordOutput{}
	err := cdc.UnmarshalJSON([]byte(body), &record)
	require.Nil(t, err)
	return record
}

func getAsset(t *testing.T, port string, assetID string) asset.RecordOutput {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", fmt.Sprintf("/assets/%s", assetID), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	record := asset.RecordOutput{}
	err := cdc.UnmarshalJSON([]byte(body), &record)
	require.Nil(t, err)
	return record
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

func getProposalByOwner(t *testing.T, port string, addr sdk.AccAddress) []asset.ProposalOutput {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", fmt.Sprintf("/accounts/%s/proposals", addr.String()), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var proposals []asset.ProposalOutput
	err := cdc.UnmarshalJSON([]byte(body), &proposals)
	require.Nil(t, err)
	return proposals
}

func getTxsProperties(t *testing.T, port string) []asset.HistoryUpdateProperty {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", fmt.Sprintf("/assets/%s/properties/size/history", "test"), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var history []asset.HistoryUpdateProperty
	err := json.Unmarshal([]byte(body), &history)
	require.Nil(t, err)
	return history
}

func getRecordOwners(t *testing.T, port string) []asset.HistoryTransferOutput {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", fmt.Sprintf("/assets/%s/owners/history", "test"), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var owners []asset.HistoryTransferOutput
	err := json.Unmarshal([]byte(body), &owners)
	require.Nil(t, err)
	return owners
}

// Test Identity Module
// -------------------------------------------------------------------------------------------------

func TestAddTrust(t *testing.T) {
	name, password := "test", "1234567890"
	addr, seed := CreateAddr(t, name, password, GetKeyBase(t))
	name2, _ := "test2", "1234567890"
	addr2, _ := CreateAddr(t, name2, password, GetKeyBase(t))
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

func TestAddCerts(t *testing.T) {
	name, password := "test", "1234567890"
	addr, seed := CreateAddr(t, name, password, GetKeyBase(t))
	cleanup, _, port := InitializeTestLCD(t, 1, []sdk.AccAddress{addr})
	defer cleanup()

	resultTx := doRegisterIdentity(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	// AddCerts tests
	resultTx = doAddCerts(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	certs := getCertsByOwner(t, port)
	assert.Equal(t, len(certs), 2)
}

func TestRegisterIdent(t *testing.T) {
	name, password := "test", "1234567890"
	addr, seed := CreateAddr(t, name, password, GetKeyBase(t))
	cleanup, _, port := InitializeTestLCD(t, 1, []sdk.AccAddress{addr})
	defer cleanup()

	// AddCerts tests
	resultTx := doRegisterIdentity(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	// Onwers
	owners := getOwners(t, port, addr)
	assert.Equal(t, len(owners), 1)

	// Test AddOwner
	resultTx = doAddOwner(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	// Onwers
	owners = getOwners(t, port, addr)
	assert.Equal(t, len(owners), 2)

	// Test DeleteOwner
	resultTx = doDelOwner(t, port, seed, name, password, addr)
	tests.WaitForHeight(resultTx.Height+1, port)
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	// Onwers
	owners = getOwners(t, port, addr)
	assert.Equal(t, len(owners), 1)

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
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
			"chain_id": "%s"
		},
		"trust": true
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", fmt.Sprintf("/idents/%s/trusts", trusting), jsonStr)
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
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
			"chain_id": "%s"
		},
		"values": [
			{
				"property": "company",
				"data": {
					"demo": "1212"
				},
				"confidence": true
			},
			{
				"id":"",
				"context":"",
				"property":"entity.person",
				"data":{  
					"address_line_1":"Hoang Liet Hoang Mai",
					"address_line_2":"",
					"city":"Hanoi",
					"corp_num":"",
					"country":"",
					"effective_date":"",
					"end_date":"",
					"first_name":"Thang",
					"last_name":"Nguyen",
					"legal_entity_id":"",
					"org_type":"",
					"postal_code":"",
					"province":"Hà Nội"
				},
				"confidence":true
			}
		]
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", fmt.Sprintf("/idents/cosmosaccaddr1753dqa50dlh8l4xl0j0kd9gga0heqsj7c2wwef/certs"), jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func getTrusts(t *testing.T, port string, trustor sdk.AccAddress) []sdk.AccAddress {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", fmt.Sprintf("/idents/%s/trusts", trustor), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var trusts []sdk.AccAddress
	err := cdc.UnmarshalJSON([]byte(body), &trusts)
	require.Nil(t, err)
	return trusts
}

func getCertsByOwner(t *testing.T, port string) []identity.Cert {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", "/idents/cosmosaccaddr1753dqa50dlh8l4xl0j0kd9gga0heqsj7c2wwef/certs", nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var certs []identity.Cert
	err := cdc.UnmarshalJSON([]byte(body), &certs)
	require.Nil(t, err)
	return certs
}

func getOwners(t *testing.T, port string, addr sdk.AccAddress) []sdk.AccAddress {
	// get the account to get the sequence
	res, body := Request(t, port, "GET", fmt.Sprintf("/idents/%s/owners", addr), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var owners []sdk.AccAddress
	err := cdc.UnmarshalJSON([]byte(body), &owners)
	require.Nil(t, err)
	return owners
}

func doRegisterIdentity(t *testing.T, port, seed, name, password string, addr sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)
	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
			"chain_id": "%s"
		}
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", fmt.Sprintf("/idents/%s/register", addr), jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func doAddOwner(t *testing.T, port, seed, name, password string, addr sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)
	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
			"chain_id": "%s"
		},
		"owner": "cosmosaccaddr1jawd35d9aq4u76sr3fjalmcqc8hqygs9gtnmv3"
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "POST", fmt.Sprintf("/idents/%s/owners", addr), jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func doDelOwner(t *testing.T, port, seed, name, password string, addr sdk.AccAddress) (resultTx ctypes.ResultBroadcastTxCommit) {
	acc := getAccount(t, port, addr)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()
	chainID := viper.GetString(client.FlagChainID)
	// send
	jsonStr := []byte(fmt.Sprintf(`{
		"base_req": {
			"name": "%s",
			"password": "%s",
			"account_number": "%d",
			"sequence": "%d",
			"gas": "10000",
			"chain_id": "%s"
		}, 
		"owner": "cosmosaccaddr1753dqa50dlh8l4xl0j0kd9gga0heqsj7c2wwef"
	}`, name, password, accnum, sequence, chainID))

	res, body := Request(t, port, "DELETE", fmt.Sprintf("/idents/%s/owners/cosmosaccaddr1jawd35d9aq4u76sr3fjalmcqc8hqygs9gtnmv3", addr), jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)
	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func TestAppLogin(t *testing.T) {
	name, password := "test", "1234567890"
	addr, _ := CreateAddr(t, name, password, GetKeyBase(t))
	cleanup, _, port := InitializeTestLCD(t, 1, []sdk.AccAddress{addr})
	defer cleanup()
	doAppSignAndVerify(t, port)

}

func doAppSignAndVerify(t *testing.T, port string) {
	jsonStr := []byte(fmt.Sprintf(`{
		"name": "test",
		"password": "1234567890",
		"nonce": "7_tKYK2eRacnzPZDHDm7jxqLMRxFTPZ5KKZ_ZXdvuLU="
	}`))

	res, body := Request(t, port, "POST", "/apps/sign", jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	b, _ := base64.StdEncoding.DecodeString(body)

	res, body = Request(t, port, "POST", "/apps/verify", b)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
}
