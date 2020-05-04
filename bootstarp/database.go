package bootstarp

import (
		_ "github.com/WebGameLinux/cms/models"
		utils "github.com/WebGameLinux/cms/utils/beego"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/orm"
		_ "github.com/go-sql-driver/mysql"
)

// 加载 数据模型 与 orm
func bootDatabase() {
		databaseInit()
}

//调用方式
//databaseInit() 或 databaseInit("w") 或 databaseInit("default") //初始化主库
//databaseInit("w","r")	//同时初始化主库和从库
//databaseInit("w")
func databaseInit(aliases ...string) {
		//如果是开发模式，则显示命令信息
		isDev := utils.IsThatRunMode("dev")
		if isDev {
				orm.Debug = isDev
		}
		if len(aliases) <= 0 {
				registerDatabases("w")
				utils.Onerror(orm.RunSyncdb("default", false, isDev))
				return
		}
		for _, alias := range aliases {
				registerDatabases(alias)
				//主库 自动建表
				if "w" == alias {
						utils.Onerror(orm.RunSyncdb("default", false, isDev))
				}
		}
}

// 注册数据库 orm
func registerDatabases(alias string) {
		if len(alias) == 0 {
				return
		}
		//连接名称
		dbAlias := alias
		if "w" == alias || "default" == alias {
				alias = "w"
				dbAlias = "default"
		}
		driver := beego.AppConfig.String("database")
		switch driver {
		case "mysql":
				registerMysql(dbAlias, alias)
		}
}

// mysql 链接解析
func registerMysql(alias string, scope string) {
		//数据库端口
		dbPort := utils.GetKvString("db_" + scope + "_port")
		//数据库IP（域名）
		dbHost := utils.GetKvString("db_" + scope + "_host")
		//数据库连接用户名
		dbPwd := utils.GetKvString("db_" + scope + "_password")
		//数据库名称
		dbName := utils.GetKvString("db_" + scope + "_database")
		//数据库连接用户名
		dbUser := utils.GetKvString("db_" + scope + "_username")
		// 链接
		url := dbUser + ":" + dbPwd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8&loc=Local"
		utils.Onerror(orm.RegisterDataBase(alias, "mysql", url, 30))
}
