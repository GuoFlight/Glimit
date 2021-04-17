package main

import (
	"dlimit/conf"
	"dlimit/handleFunction"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

/*
函数作用：往指定目录里追加字符串，如果已存在就不追加
参数：
	file：文件名
	s：需要追加的字符串
	regular：应该传入一个正则表达式。用于检测文件是否存在此字符串
*/
func append_string(file,s,regular string){
	data,err := ioutil.ReadFile(file)
	if err!=nil {
		fmt.Printf("%s打开失败\n",file)
		os.Exit(1)
	}
	isMatch, _ := regexp.MatchString(regular, string(data))	//如果已存在就不用再添加了
	if !isMatch{
		f,err := os.OpenFile(file,os.O_RDWR |os.O_APPEND,0644)
		defer f.Close()
		if err!=nil{
			fmt.Printf("%s打开失败\n",file)
			os.Exit(1)
		}
		_,err=f.Write([]byte(s))
		if err!=nil{
			fmt.Printf("%s写入失败\n",file)
			os.Exit(1)
		}
	}
}

/*
函数作用：用户登录时自动发送请求
*/
func append_curl(){
	dlimit_bashrc :="function dlimit(){\n\tcur_user=`whoami`\n\tpid=$$\n\tcurl -s http://127.0.0.1:1216/dlimit/v1/limit_user?username=${cur_user}\\&pid=${pid}\n}\ndlimit >> /dev/null\n"
	append_string("/etc/bashrc",dlimit_bashrc,"\ndlimit")
}

/*
函数作用：添加计划任务，执行清理脚本，自动清理过期cgroup组
*/
func append_crontab(){
	dir_clean,err := filepath.Abs(".")	//得到当前程序的父目录
	if err!=nil{
		fmt.Println("无法得到程序的父目录")
		return
	}
	path_clean := fmt.Sprintf("%s/clean_cgroup.sh",dir_clean)	//得到当前程序的路径


	//替换脚本中的cgroup目录
	data_before,err := ioutil.ReadFile(path_clean)		//替换前的脚本
	if err!=nil{
		fmt.Println("清理脚本无法读取，请检查")
		return
	}
	data_after := strings.Replace(string(data_before),"{CGROUP_DIR}",conf.Config.Cgroup.Dir,1)	//替换后的脚本
	ioutil.WriteFile(path_clean,[]byte(data_after),0755)

	//给脚本添加执行权限
	cmd := exec.Command("chmod", "+x", fmt.Sprintf("%s/clean_cgroup.sh",dir_clean))
	err = cmd.Run()
	if err!=nil{
		fmt.Printf("请检查%s/clean_cgroup.sh是否有权限\n",dir_clean)
	}

	//将计划任务写入/etc/crontab
	job := fmt.Sprintf("0 4 * * * root %s/clean_cgroup.sh",dir_clean)
	job_regular := strings.Replace("\n"+job,"*","\\*",3)
	append_string("/etc/crontab",job+"\n",job_regular)
}

func main() {
	//加载配置文件
	conf.Load_conf()
	fmt.Println(conf.Config)

	//用户登录时自动发送请求
	append_curl()

	//添加计划任务，执行清理脚本，自动清理过期cgroup组
	append_crontab()

	//路由
	router := gin.Default()         //Default 使用 Logger 和 Recovery 中间件
	router.GET("/dlimit/v1/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/dlimit/v1/limit_user", handleFunction.Handle_limit)
	router.Run(":"+strconv.Itoa(conf.Config.Port))        //启动服务
}
