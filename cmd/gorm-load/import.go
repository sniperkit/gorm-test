package main

import (
	"runtime"
)

func loadFileDB(filePath string) {
	wrap := "\n"
	if runtime.GOOS == "windows" {
		wrap = "\r\n"
	}
	mysql.RegisterLocalFile(filePath)
	loadSql := "load data LOW_PRIORITY local infile '%s' into table url FIELDS TERMINATED BY ','  LINES TERMINATED BY '%s' (Url,`Name`)  set UrlGroupId=1"
	loadSql = fmt.Sprintf(loadSql, filePath, wrap)
	_, err := dbs.Exec(loadSql)
	if err != nil {
		glog.Error("sql is warn, sql: %s err: %s \n", loadSql, err.Error())
		sqlClose()
		sqlConntion()
		return
	}
	glog.Info("db run success! filePath: %s \n", filePath)
}
