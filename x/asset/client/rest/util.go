package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/icheckteam/ichain/client/errors"
	"github.com/icheckteam/ichain/x/asset"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

func signAndBuild(ctx context.CoreContext, cdc *wire.Codec, w http.ResponseWriter, m baseBody, msg sdk.Msg) {
	ctx = ctx.WithGas(m.Gas)
	ctx = ctx.WithAccountNumber(m.AccountNumber)
	ctx = ctx.WithSequence(m.Sequence)
	ctx = ctx.WithChainID(m.ChainID)

	if len(m.Memo) > 0 {
		ctx = ctx.WithMemo(m.Memo)
	}

	txBytes, err := ctx.SignAndBuild(m.Name, m.Password, []sdk.Msg{msg}, cdc)

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

func withErrHandler(fn func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			errors.WriteError(w, err)
			return
		}
	}
}

func getReporters(ctx context.CoreContext, recordID string, cdc *wire.Codec) ([]asset.Reporter, error) {
	reportersPrefixKey := asset.GetReportersKey(recordID)
	kvs, err := ctx.QuerySubspace(cdc, reportersPrefixKey, storeName)
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

func getProperties(ctx context.CoreContext, recordID string, cdc *wire.Codec) ([]asset.Property, error) {
	propertiesPrefixKey := asset.GetPropertiesKey(recordID)
	kvs, err := ctx.QuerySubspace(cdc, propertiesPrefixKey, storeName)
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

func getMaterials(ctx context.CoreContext, recordID string, cdc *wire.Codec) ([]asset.Material, error) {
	materialsPrefixKey := asset.GetMaterialsKey(recordID)
	kvs, err := ctx.QuerySubspace(cdc, materialsPrefixKey, storeName)
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

func widthMoreRecord(ctx context.CoreContext, record asset.Asset, cdc *wire.Codec, includes ...string) (*asset.RecordOutput, error) {
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
		Unit:     record.Unit,
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

func getRecord(ctx context.CoreContext, recordID string, cdc *wire.Codec, includes ...string) (*asset.RecordOutput, error) {
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
	if recordOutput.Root != "" {
		root, err := getRecord(ctx, recordOutput.Root, cdc, "properties")
		if err != nil {
			return nil, err
		}
		record.Unit = root.Unit
		for _, p := range root.Properties {
			switch p.Name {
			case "barcode":
				recordOutput.Barcode = p.StringValue
				break
			case "type":
				recordOutput.Type = p.StringValue
				break
			case "subtype":
				recordOutput.SubType = p.StringValue
				break
			default:
				break
			}
		}
	}

	return recordOutput, nil
}

func getRecordsByAccount(ctx context.CoreContext, addr sdk.AccAddress, cdc *wire.Codec) ([]*asset.RecordOutput, error) {
	recordsPrefixKey := asset.GetAccountAssetsKey(addr)
	kvs, err := ctx.QuerySubspace(cdc, recordsPrefixKey, storeName)
	if err != nil {
		return nil, err
	}
	return getRecordsByKvs(ctx, kvs, cdc)
}

func getRecordsByKvs(ctx context.CoreContext, kvs []sdk.KVPair, cdc *wire.Codec) ([]*asset.RecordOutput, error) {
	records := make([]*asset.RecordOutput, len(kvs))
	for i, kv := range kvs {
		var recordID string
		err := cdc.UnmarshalBinary(kv.Value, &recordID)
		if err != nil {
			return nil, fmt.Errorf("getRecordsByKvs %s err: %s", recordID, err.Error())
		}
		record, err := getRecord(ctx, recordID, cdc)
		if err != nil {
			return nil, err
		}
		records[i] = record
	}
	return records, nil
}

func getProposals(ctx context.CoreContext, kvs []sdk.KVPair, cdc *wire.Codec) (asset.Proposals, error) {
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
