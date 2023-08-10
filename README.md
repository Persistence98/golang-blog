# golang-blog
该项目为golang基础框架。

基础框架：gin、gorm

数据库：mysql、redis

内置基础：jwt鉴权、限流器（令牌桶）等


## 下载使用
第一步 
git clone https://github.com/Persistence98/golang-blog.git

第二步
配置config.yaml

<div class="code-block">
```mysql: root:password@tcp(127.0.0.1:3306)/databases_name?charset=utf8&parseTime=true&loc=Local```
<button class="copy-button" data-clipboard-target=".code-block">复制</button>
</div>


第三步
使用命令下载gin、gorm、jwt等第三方扩展
go mod tidy

第四步
go run main.go
注：必须在main.go同级目录
