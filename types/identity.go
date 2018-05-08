package types

// Identity ....
type Identity struct {
	ClaimID string
	Name    string
	Lock    bool
}

// Identities ...
type Identities []*Identity
