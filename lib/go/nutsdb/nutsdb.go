//go:build wasip1

package nutsdb

import (
	"context"

	gnutsdb "github.com/iansmith/parigot/g/nutsdb/v1"
)

func WritePair(ctx context.Context, client gnutsdb.Client, id gnutsdb.NutsDBId, bucketPath, key string, value []byte) *gnutsdb.FutureWritePair {
	pair := &gnutsdb.Pair{
		BucketPath: bucketPath,
		Key:        []byte(key),
		Value:      []byte(value),
	}
	req := &gnutsdb.WritePairRequest{NutsdbId: id.Marshal(), Pair: pair}

	return client.WritePair(ctx, req)
}

func ReadPair(ctx context.Context, client gnutsdb.Client, id gnutsdb.NutsDBId, bucketPath, key string, defaultValue []byte) *gnutsdb.FutureReadPair {
	pair := &gnutsdb.Pair{
		BucketPath: bucketPath,
		Key:        []byte(key),
		Value:      []byte(defaultValue),
	}
	req := &gnutsdb.ReadPairRequest{NutsdbId: id.Marshal(), Pair: pair}

	return client.ReadPair(ctx, req)
}

func Open(ctx context.Context, dbname string) *gnutsdb.FutureOpen {
	req := &gnutsdb.OpenRequest{}
	req.DbName = dbname
	return gnutsdb.OpenHost(ctx, req)
}

func Close(ctx context.Context, client gnutsdb.Client, id gnutsdb.NutsDBId) *gnutsdb.FutureClose {
	req := &gnutsdb.CloseRequest{}
	req.NutsdbId = id.Marshal()
	return client.Close(ctx, req)
}
