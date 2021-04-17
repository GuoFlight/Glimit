#!/bin/bash

############################################
# 作者：郭飞
# 作用：清理过期的cgroup组
# 注意：应结合计划任务一起使用
############################################

subsystems=("cpu" "memory")		#需要清理的子系统
dir_cgroup="{CGROUP_DIR}"	    #cgroup目录
tody=`date "+%Y%m%d"`

for i in ${subsystems[*]}; do
    groups=`ls $dir_cgroup/$i | grep itools`
    for j in $groups; do
        date=`echo $j | awk -F'_' '{print $2}'`
        [ ! -z $date ] && [ $date -lt $tody ] && cgdelete $i:$j
    done
done





