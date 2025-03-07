package bluesky

import (
	"context"
	"fmt"
	"time"

	"github.com/bluesky-social/indigo/xrpc"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func ExpiresAt(ctx context.Context, auth *xrpc.AuthInfo) (time.Time, error) {
	// pubKey, err := GetPublicKey(ctx, auth.Did)
	// if err != nil {
	// 	return time.Time{}, fmt.Errorf("error getting public key: %w", err)
	// }

	// parsedPubKey, err := btcec.ParsePubKey(pubKey.Bytes())
	// if err != nil {
	// 	return time.Time{}, fmt.Errorf("failed to parse public key: %w", err)
	// }

	// slog.Debug("parsed public key", slog.String("x", parsedPubKey.ToECDSA().X.String()), slog.String("y", parsedPubKey.ToECDSA().Y.String()), slog.Bool("curve", parsedPubKey.ToECDSA() == elliptic.P256()), slog.String("curve_name", parsedPubKey.ToECDSA().Params().Name))

	// // https://atproto.com/specs/oauth#demonstrating-proof-of-possession-d-po-p for ES256 (actually ES256K)
	// verifiedToken, err := jwt.Parse([]byte(auth.AccessJwt), jwt.WithKey(jwa.ES256K(), parsedPubKey.ToECDSA()))
	// if err != nil {
	// 	return time.Time{}, fmt.Errorf("failed to verify jws: %w", err)
	// }

	token, err := jwt.ParseInsecure([]byte(auth.AccessJwt))
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse jwt: %w", err)
	}

	expr, _ := token.Expiration()

	return expr, nil

}
