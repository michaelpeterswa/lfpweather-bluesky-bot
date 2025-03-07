package bluesky

// import (
// 	"context"
// 	"fmt"

// 	"github.com/bluesky-social/indigo/atproto/crypto"
// 	"github.com/bluesky-social/indigo/atproto/identity"
// 	"github.com/bluesky-social/indigo/atproto/syntax"
// )

// func GetPublicKey(ctx context.Context, did string) (crypto.PublicKey, error) {
// 	syntaxDID, err := syntax.ParseDID(did)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse DID: %w", err)
// 	}

// 	ident, err := identity.DefaultDirectory().LookupDID(ctx, syntaxDID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to lookup DID: %w", err)
// 	}

// 	key, err := ident.PublicKey()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get public key: %w", err)
// 	}

// 	return key, nil
// }
