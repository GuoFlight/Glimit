package conf

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path/filepath"
)

const (
	conf_path = "./dlimit.toml"
)

type Config_file struct {
	Whitelist []string `toml:"whitelist"`
	Port int `toml:"port"`
	Cgroup struct{
		Dir string `toml:"dir"`
	} `toml:"cgroup"`
	Cpu struct{
		Cfs_period_us int `toml:"cfs_period_us"`
		Cfs_quota_us int  `toml:"cfs_quota_us"`
	} `toml:"cpu"`
}
var (
	Config Config_file
)

func Load_conf(){
	//加载配置文件
	if _, err := toml.DecodeFile(conf_path, &Config); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	//纠正配置文件
	var err error
	if Config.Cgroup.Dir,err = filepath.Abs(Config.Cgroup.Dir);err!=nil{
		log.Fatal("cgroup路径格式错误")
		os.Exit(1)
	}
}
