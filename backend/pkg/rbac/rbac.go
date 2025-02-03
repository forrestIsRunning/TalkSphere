package rbac

import (
	"TalkSphere/setting"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Enforcer *casbin.Enforcer

func InitCasbin() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		setting.Conf.MysqlConfig.User,
		setting.Conf.MysqlConfig.PassWord,
		setting.Conf.MysqlConfig.Host,
		setting.Conf.MysqlConfig.DB,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("连接数据库失败", zap.Error(err))
	}

	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(db, &TestCasbinRule{}, "casbin_rule")
	if err != nil {
		zap.L().Fatal("初始化 Casbin 适配器失败", zap.Error(err))
	}

	Enforcer, err = casbin.NewEnforcer("conf/rbac_model.conf", adapter)
	if err != nil {
		zap.L().Fatal("初始化 Casbin 失败", zap.Error(err))
	}

	loadDefaultPolicy()

	if err := Enforcer.LoadPolicy(); err != nil {
		zap.L().Fatal("加载 Casbin 策略失败", zap.Error(err))
	}
}

func loadDefaultPolicy() {
	_, _ = Enforcer.AddPolicy("admin", "/api/*", "GET")
	_, _ = Enforcer.AddPolicy("admin", "/api/*", "POST")
	_, _ = Enforcer.AddPolicy("admin", "/api/*", "PUT")
	_, _ = Enforcer.AddPolicy("admin", "/api/*", "DELETE")
	_, _ = Enforcer.AddPolicy("admin", "/api/*", "PATCH")

	_, _ = Enforcer.AddPolicy("user", "/api/login", "POST")
	_, _ = Enforcer.AddPolicy("user", "/api/register", "POST")

	_ = Enforcer.SavePolicy()
}

type TestCasbinRule struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Ptype     string `gorm:"size:16"`
	V0        string `gorm:"size:128"`
	V1        string `gorm:"size:128"`
	V2        string `gorm:"size:256"`
	DeletedAt gorm.DeletedAt
}
