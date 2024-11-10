package data

import (
	"comment/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	hashid "github.com/xinghe903/xinghe/pkg/distribute/id"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewCommentRepo, NewSnowflake)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}

func NewSnowflake(c *conf.Server) *hashid.Snowflake {
	return hashid.NewSnowflake(c.DnsAddr)
}
