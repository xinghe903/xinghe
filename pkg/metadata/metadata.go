package metadata

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"
)

const (
	UserIdKey   = "user_id"
	UserNameKey = "user_name"
)

var (
	ErrNotFoundMetadata = errors.New("not found metadata")
)

func getMetadataFromKey(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrNotFoundMetadata
	}
	return md[key][0], nil
}

func setMetadataFromKey(ctx context.Context, key string, val ...string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	md[key] = val
	return metadata.NewOutgoingContext(ctx, md)
}

func GetUserId(ctx context.Context) (string, error) {
	return getMetadataFromKey(ctx, UserIdKey)
}

func GetUserName(ctx context.Context) (string, error) {
	return getMetadataFromKey(ctx, UserNameKey)
}

func SetUserId(ctx context.Context, userId string) context.Context {
	return setMetadataFromKey(ctx, UserIdKey, userId)
}

func SetUserName(ctx context.Context, userName string) context.Context {
	return setMetadataFromKey(ctx, UserNameKey, userName)
}
