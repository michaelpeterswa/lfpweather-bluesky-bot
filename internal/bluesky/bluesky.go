package bluesky

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
)

type BlueskyClient struct {
	xrpcClient           *xrpc.Client
	did                  string
	mu                   sync.Mutex
	accessTokenExpiresAt time.Time
	forceRefreshDuration time.Duration
}

func NewBlueskyClient(ctx context.Context, host string, username string, password string, forceRefreshDuration time.Duration) (*BlueskyClient, error) {
	xrpcClient := &xrpc.Client{
		Host: host,
	}

	handle, err := atproto.IdentityResolveHandle(ctx, xrpcClient, username)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve handle: %w", err)
	}

	auth, err := atproto.ServerCreateSession(ctx, xrpcClient, &atproto.ServerCreateSession_Input{
		Identifier: handle.Did,
		Password:   password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	xrpcClient.Auth = &xrpc.AuthInfo{
		AccessJwt:  auth.AccessJwt,
		RefreshJwt: auth.RefreshJwt,
		Handle:     auth.Handle,
		Did:        auth.Did,
	}

	accessTokenExpires, err := ExpiresAt(ctx, xrpcClient.Auth)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token expiration: %w", err)
	}

	return &BlueskyClient{
		xrpcClient:           xrpcClient,
		did:                  handle.Did,
		mu:                   sync.Mutex{},
		accessTokenExpiresAt: accessTokenExpires,
		forceRefreshDuration: forceRefreshDuration,
	}, nil
}

func (bc *BlueskyClient) RefreshAuth(ctx context.Context) error {
	if time.Now().Before(bc.accessTokenExpiresAt.Add(-bc.forceRefreshDuration)) {
		expiresIn := time.Until(bc.accessTokenExpiresAt)
		slog.Debug("skipping refresh auth", slog.String("did", bc.did), slog.String("expires_in", expiresIn.String()))
		return nil
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.xrpcClient.Auth.AccessJwt = bc.xrpcClient.Auth.RefreshJwt
	resp, err := atproto.ServerRefreshSession(ctx, bc.xrpcClient)
	if err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	bc.xrpcClient.Auth = &xrpc.AuthInfo{
		AccessJwt:  resp.AccessJwt,
		RefreshJwt: resp.RefreshJwt,
		Handle:     resp.Handle,
		Did:        resp.Did,
	}

	// make sure to update the expiration time
	bc.accessTokenExpiresAt, err = ExpiresAt(ctx, bc.xrpcClient.Auth)
	if err != nil {
		return fmt.Errorf("failed to get access token expiration: %w", err)
	}

	slog.Debug("refreshed auth", slog.String("did", resp.Did), slog.String("handle", resp.Handle))

	return nil
}

func (bc *BlueskyClient) WritePost(ctx context.Context, text string) error {
	currentTime := time.Now()
	post := &bsky.FeedPost{
		Text:      text,
		CreatedAt: currentTime.Format(time.RFC3339),
		Langs:     []string{"en-US"},
	}

	_, err := atproto.RepoCreateRecord(
		ctx,
		bc.xrpcClient,
		&atproto.RepoCreateRecord_Input{
			Repo:       bc.did,
			Collection: "app.bsky.feed.post",
			Record:     &util.LexiconTypeDecoder{Val: post},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create record: %w", err)
	}

	return nil
}
