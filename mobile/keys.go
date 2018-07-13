package mobile

import (
	"github.com/icheckteam/ichain/crypto/keys"
	tcrypto "github.com/tendermint/tendermint/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
)

type KeyStore struct {
	keystore keys.Keybase
}

// NewKeyStore creates a keystore for the given directory.
func NewKeyStore(keydir string) *KeyStore {
	return &KeyStore{keystore: keys.New(dbm.NewFSDB(keydir))}
}

// List returns the keys from storage in alphabetical order.
func (ks KeyStore) List() ([]keys.Info, error) {
	return ks.keystore.List()
}

// Get returns the public information about one key.
func (ks KeyStore) Get(name string) (keys.Info, error) {
	return ks.keystore.Get(name)
}

// CreateMnemonic generates a new key and persists it to storage, encrypted
// using the passphrase.  It returns the generated seedphrase
// (mnemonic) and the key Info.  It returns an error if it fails to
// generate a key for the given algo type, or if another key is
// already stored under the same name.
func (ks KeyStore) CreateMnemonic(name string, language keys.Language, passwd string, algo keys.SigningAlgo) (keys.Info, string, error) {
	return ks.keystore.CreateMnemonic(name, language, passwd, algo)
}

// TEMPORARY METHOD UNTIL WE FIGURE OUT USER FACING HD DERIVATION API
func (kb KeyStore) CreateKey(name, mnemonic, passwd string) (keys.Info, error) {
	return kb.keystore.CreateKey(name, mnemonic, passwd)
}

// Sign signs the msg with the named key.
// It returns an error if the key doesn't exist or the decryption fails.
func (kb KeyStore) Sign(name, passphrase string, msg []byte) (sig tcrypto.Signature, pub tcrypto.PubKey, err error) {
	return kb.keystore.Sign(name, passphrase, msg)
}

func (kb KeyStore) Export(name string) (armor string, err error) {
	return kb.keystore.Export(name)
}

// ExportPubKey returns public keys in ASCII armored format.
// Retrieve a Info object by its name and return the public key in
// a portable format.
func (kb KeyStore) ExportPubKey(name string) (armor string, err error) {
	return kb.keystore.ExportPubKey(name)
}

func (kb KeyStore) Import(name string, armor string) (err error) {
	return kb.keystore.Import(name, armor)
}

// ImportPubKey imports ASCII-armored public keys.
// Store a new Info object holding a public key only, i.e. it will
// not be possible to sign with it as it lacks the secret key.
func (kb KeyStore) ImportPubKey(name string, armor string) (err error) {
	return kb.keystore.Import(name, armor)
}

// Delete removes key forever, but we must present the
// proper passphrase before deleting it (for security).
// A passphrase of 'yes' is used to delete stored
// references to offline and Ledger / HW wallet keys
func (kb KeyStore) Delete(name, passphrase string) error {
	return kb.keystore.Delete(name, passphrase)
}

// Update changes the passphrase with which an already stored key is
// encrypted.
//
// oldpass must be the current passphrase used for encryption,
// getNewpass is a function to get the passphrase to permanently replace
// the current passphrase
func (kb KeyStore) Update(name, oldpass, newPass string) error {
	return kb.keystore.Update(name, oldpass, newPass)

}
