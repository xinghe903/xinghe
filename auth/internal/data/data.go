package data

import (
	"auth/internal/conf"
	"fmt"

	"auth/internal/data/po"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	hashid "github.com/xinghe903/xinghe/pkg/distribute/id"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewSnowflake, NewGormClient, NewUserRepo,
	NewPermissionRepo, NewRolePermissionRepo, NewRoleRepo, NewAuthRepo, NewUserRoleRepo,
)

// Data .
type Data struct {
	db  *gorm.DB
	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	data := &Data{
		db:  db,
		log: log.NewHelper(logger),
	}
	cleanup := func() {
		data.log.Info("closing the data resources")
	}
	err := checkTable(db, data.log)
	if err != nil {
		return nil, nil, err
	}
	return data, cleanup, nil
}

func checkTable(dbIns *gorm.DB, log *log.Helper) error {
	entries := []interface{}{&po.User{}, &po.Auth{}, &po.Role{}, &po.RolePermission{}, &po.Permission{}}
	for _, it := range entries {
		action := "创建"
		if dbIns.Migrator().HasTable(it) {
			action = "更新"
			log.Infof("表%s: %v", action, it)
		} else {
			log.Infof("表%s: %v", action, it)
		}
		if err := dbIns.AutoMigrate(it); err != nil {
			return err
		}
		if dbIns.Migrator().HasTable(it) {
			log.Infof("表%s成功", action)
		} else {
			log.Infof("表%s失败", action)
		}
	}
	return nil
}

func NewGormClient(conf *conf.Data) (*gorm.DB, error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Database,
		true,
		"Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(int(conf.Database.MaxOpenConnections))
	sqlDB.SetConnMaxLifetime(conf.Database.MaxConnectionLifeTime.AsDuration())
	sqlDB.SetMaxIdleConns(int(conf.Database.MaxIdleConnections))
	return db, nil
}

func NewSnowflake() *hashid.Sonyflake {
	return hashid.NewSonyflake()
}
