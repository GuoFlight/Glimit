
all:
	@echo "start building..."
	go build -mod=vendor -o dlimit-server
	@rm -rf dlimit
	@mkdir dlimit
	@cp -a dlimit-server dlimit/
	@cp -a dlimit.toml dlimit/
	@chmod +x clean_cgroup.sh && cp -a clean_cgroup.sh dlimit/
	@echo "编译完成，请将dlimit目录移动到适当的目录，并启动：sudo nohup ./dlimit-server >> /dev/null &"
