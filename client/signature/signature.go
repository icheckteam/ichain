package signature

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/crypto"
)

// LoginBody ...
type LoginBody struct {
	Website  string `json:"webiste"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Nonce    string `json:"nonce"`
}

type ClaimMsg struct {
	PubKey  crypto.PubKey `json:"pubkey"`
	Expires int64         `json:"expires"`
	Nonce   string        `json:"nonce"`
}

func (msg ClaimMsg) Bytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

type Claim struct {
	Msg       ClaimMsg         `json:"msg"`
	Signature crypto.Signature `json:"signature"`
}

///////////////////////////
// REST

// get key REST handler
func SignHandler(w http.ResponseWriter, r *http.Request) {
	var m LoginBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	keybase, err := keys.GetKeyBase()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	key, err := keybase.Get(m.Name)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	msg := ClaimMsg{
		Nonce:   m.Nonce,
		PubKey:  key.GetPubKey(),
		Expires: time.Now().Add(60 * time.Second).Unix(),
	}

	sign, _, err := keybase.Sign(m.Name, m.Password, msg.Bytes())
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	output, err := cdc.MarshalJSON(&Claim{
		Msg:       msg,
		Signature: sign,
	})
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(output)
}

type VerifyOutput struct {
	Status bool
}

func VerifiyHandler(w http.ResponseWriter, r *http.Request) {
	var m Claim
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	err = cdc.UnmarshalJSON(body, &m)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	output := VerifyOutput{
		Status: verify(m),
	}

	b, err := json.Marshal(output)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(b)

}

func verify(m Claim) bool {
	if m.Msg.Expires < time.Now().Unix() {
		return false
	}
	return m.Msg.PubKey.VerifyBytes(m.Msg.Bytes(), m.Signature)
}

// resgister REST routes
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/apps/sign", SignHandler).Methods("POST")
	r.HandleFunc("/apps/verify", VerifiyHandler).Methods("POST")
}
