package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/icheckteam/ichain/client/errors"
	"github.com/icheckteam/ichain/x/asset"
)

type bodyI interface {
	ValidateBasic() error
}

func signAndBuild(ctx context.CLIContext, cdc *wire.Codec, w http.ResponseWriter, m baseBody, msg sdk.Msg) {

	txCtx := authctx.TxContext{
		Codec:         cdc,
		Gas:           m.Gas,
		ChainID:       m.ChainID,
		AccountNumber: m.AccountNumber,
		Sequence:      m.Sequence,
		Memo:          m.Memo,
	}

	txBytes, err := txCtx.BuildAndSign(m.Name, m.Password, []sdk.Msg{msg})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	// send
	res, err := ctx.BroadcastTx(txBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BroadcastTx:" + err.Error()))
		return
	}

	output, err := wire.MarshalJSONIndent(cdc, res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(output)
}

// WriteJSON ...
func WriteJSON(w http.ResponseWriter, cdc *wire.Codec, data interface{}) {
	output, err := cdc.MarshalJSON(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(output)
}

// WriteJSON2 ...
func WriteJSON2(w http.ResponseWriter, cdc *wire.Codec, data interface{}) {
	output, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(output)
}

func validateAndGetDecodeBody(r *http.Request, cdc *wire.Codec, m bodyI) error {
	body, err := ioutil.ReadAll(r.Body)
	err = cdc.UnmarshalJSON(body, m)
	if err != nil {
		return err
	}

	if err := m.ValidateBasic(); err != nil {
		return err
	}
	return nil
}

func withErrHandler(fn func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			errors.WriteError(w, err)
			return
		}
	}
}

func getReporters(ctx context.CLIContext, recordID string, cdc *wire.Codec) ([]asset.Reporter, error) {
	reportersPrefixKey := asset.GetReportersKey(recordID)
	kvs, err := ctx.QuerySubspace(reportersPrefixKey, storeName)
	if err != nil {
		return nil, err
	}
	reporters := make([]asset.Reporter, len(kvs))
	for i, kv := range kvs {
		reporter, err := asset.UnmarshalReporter(cdc, kv.Value)
		if err != nil {
			return nil, err
		}
		reporters[i] = reporter
	}
	return reporters, nil
}

func getProperties(ctx context.CLIContext, recordID string, cdc *wire.Codec) ([]asset.Property, error) {
	propertiesPrefixKey := asset.GetPropertiesKey(recordID)
	kvs, err := ctx.QuerySubspace(propertiesPrefixKey, storeName)
	if err != nil {
		return nil, err
	}
	properties := make([]asset.Property, len(kvs))
	for i, kv := range kvs {
		property, err := asset.UnmarshalProperty(cdc, kv.Value)
		if err != nil {
			return nil, err
		}
		properties[i] = property
	}
	return properties, nil
}

func getMaterials(ctx context.CLIContext, recordID string, cdc *wire.Codec) ([]asset.Material, error) {
	materialsPrefixKey := asset.GetMaterialsKey(recordID)
	kvs, err := ctx.QuerySubspace(materialsPrefixKey, storeName)
	if err != nil {
		return nil, err
	}
	materials := make([]asset.Material, len(kvs))
	for i, kv := range kvs {
		material, err := asset.UnmarshalMaterial(cdc, kv.Value)
		if err != nil {
			return nil, err
		}
		materials[i] = material
	}
	return materials, nil
}

func widthMoreRecord(ctx context.CLIContext, record asset.Asset, cdc *wire.Codec, includes ...string) (*asset.RecordOutput, error) {
	recordOutput := asset.RecordOutput{
		ID:       record.ID,
		Name:     record.Name,
		Owner:    record.Owner,
		Parent:   record.Parent,
		Root:     record.Root,
		Final:    record.Final,
		Quantity: record.Quantity,
		Height:   record.Height,
		Created:  record.Created,
	}

	// defaults
	if len(includes) == 0 {
		includes = []string{
			"properties", "materials", "reporters",
		}
	}

	for _, include := range includes {
		switch include {
		case "properties":
			// query all properties of this record
			properties, err := getProperties(ctx, record.ID, cdc)
			if err != nil {
				return nil, err
			}
			recordOutput.Properties = properties
			break
		case "materials":
			// query all materials of this record
			materials, err := getMaterials(ctx, record.ID, cdc)
			if err != nil {
				return nil, err
			}
			recordOutput.Materials = materials
			break
		case "reporters":
			// query all reporters of this record
			reporters, err := getReporters(ctx, record.ID, cdc)
			if err != nil {
				return nil, err
			}
			recordOutput.Reporters = reporters
			break
		}
	}
	return &recordOutput, nil
}

func getRecord(ctx context.CLIContext, recordID string, cdc *wire.Codec, includes ...string) (*asset.RecordOutput, error) {
	recordKey := asset.GetAssetKey(recordID)
	res, err := ctx.QueryStore(recordKey, storeName)
	if err != nil {
		return nil, err
	}

	record, err := asset.UnmarshalRecord(cdc, res)
	if err != nil {
		return nil, err
	}

	recordOutput, err := widthMoreRecord(ctx, record, cdc, includes...)
	if err != nil {
		return nil, err
	}
	// get rppot asset info
	props := recordOutput.Properties
	if recordOutput.Root != "" {
		root, err := getRecord(ctx, recordOutput.Root, cdc, "properties")
		if err != nil {
			return nil, err
		}
		props = root.Properties
	}
	formatRecordProperties(recordOutput, props)
	return recordOutput, nil
}

func formatRecordProperties(record *asset.RecordOutput, props asset.Properties) {
	for _, p := range props {
		switch p.Name {
		case "barcode":
			record.Barcode = p.StringValue
			break
		case "unit":
			record.Barcode = p.StringValue
			break
		case "type":
			record.Type = p.StringValue
			break
		case "subtype":
			record.SubType = p.StringValue
			break
		default:
			break
		}
	}
}

func getRecordsByAccount(ctx context.CLIContext, addr sdk.AccAddress, cdc *wire.Codec) (asset.RecordsOutput, error) {
	recordsPrefixKey := asset.GetAccountAssetsKey(addr)
	kvs, err := ctx.QuerySubspace(recordsPrefixKey, storeName)
	if err != nil {
		return nil, err
	}
	return getRecordsByKvs(ctx, kvs, cdc)
}

func getRecordsByKvs(ctx context.CLIContext, kvs []sdk.KVPair, cdc *wire.Codec) (asset.RecordsOutput, error) {
	records := make(asset.RecordsOutput, len(kvs))
	for i, kv := range kvs {
		recordID := string(kv.Key[1+sdk.AddrLen:])
		record, err := getRecord(ctx, recordID, cdc)
		if err != nil {
			return nil, err
		}
		records[i] = *record
	}
	return records, nil
}

func getProposals(ctx context.CLIContext, kvs []sdk.KVPair, cdc *wire.Codec) (asset.Proposals, error) {
	proposals := make(asset.Proposals, len(kvs))
	for i, kv := range kvs {
		proposal, err := asset.UnmarshalProposal(cdc, kv.Value)
		if err != nil {
			return nil, err
		}
		proposals[i] = proposal
	}
	return proposals, nil
}

func getProposal(ctx context.CLIContext, addr sdk.AccAddress, recordID string, cdc *wire.Codec) (proposal asset.Proposal, err error) {
	res, err := ctx.QueryStore(asset.GetProposalKey(recordID, addr), storeName)
	if err != nil {
		return
	}
	proposal, err = asset.UnmarshalProposal(cdc, res)
	return
}
