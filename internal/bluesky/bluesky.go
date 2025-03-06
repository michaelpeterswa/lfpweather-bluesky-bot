package bluesky

import (
	"context"
	"fmt"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
)

type BlueskyClient struct {
	xrpcClient *xrpc.Client
	did        string
}

func NewBlueskyClient(ctx context.Context, host string, username string, password string) (*BlueskyClient, error) {
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

	return &BlueskyClient{
		xrpcClient: xrpcClient,
		did:        handle.Did,
	}, nil
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
