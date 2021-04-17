package handleFunction

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/exec"
	"strconv"
	"time"
	"dlimit/conf"
)
func handleErr(c *gin.Context,msg string){
	c.JSON(200, gin.H{
		"status": "error",
		"message":msg,
	})
}

func Handle_limit(c *gin.Context) {
	tody := time.Now().Format("20060102")
	username := c.Query("username")		//用户名
	pid_string := c.Query("pid")				//进程id


	//检测参数是否正确
	if username == "" || pid_string == ""{
		handleErr(c,"参数错误，请指定正确的username或pid")
		return
	}
	_,err:=strconv.Atoi(pid_string)
	if err != nil{	//转换失败则错误
		handleErr(c,"参数错误，请指定正确的pid")
		return
	}

	//写入cgroup
	isExist := false	//查看此用户是否在白名单去
	for _,v:= range conf.Config.Whitelist{
		if username == v{
			isExist = true
		}
	}
	if isExist{		//如果此用户在白名单内
		c.JSON(200, gin.H{
			"status": "success",
			"message":"此用户将不受限制",
		})
		return
	}
	dir := fmt.Sprintf("%s/cpu/itools_%s_%s", conf.Config.Cgroup.Dir, tody, username)
	if _,err:=os.Stat(dir); err != nil && os.IsNotExist(err){	//判断是否不存在
		cmd := exec.Command("cgcreate", "-g", fmt.Sprintf("cpu:itools_%s_%s",tody,username))		//创建cgrouop组
		err := cmd.Run()
		cmd = exec.Command("cgset", "-r","cpu.cfs_period_us="+strconv.Itoa(conf.Config.Cpu.Cfs_period_us), fmt.Sprintf("itools_%s_%s",tody,username))
		err = cmd.Run()
		cmd = exec.Command("cgset", "-r","cpu.cfs_quota_us="+strconv.Itoa(conf.Config.Cpu.Cfs_quota_us), fmt.Sprintf("itools_%s_%s",tody,username))
		err = cmd.Run()
		if err != nil{
			handleErr(c,"cgroup操作失败")
			return
		}
	}
	//将进程号写入cgroup
	f,err := os.OpenFile(fmt.Sprintf("%s/cgroup.procs",dir),os.O_RDWR |os.O_APPEND,0777)
	defer f.Close()
	if err!=nil{
		handleErr(c,fmt.Sprintf("%s/cgroup.procs文件打开失败",dir))
		return
	}
	_,err=f.Write([]byte(pid_string))
	if err!=nil{
		handleErr(c,"cgroup写入失败")
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"message":"",
	})
	return
}
