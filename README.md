[![Build Status](https://travis-ci.com/DaoYoung/gen-model.svg?branch=master)](https://travis-ci.com/DaoYoung/gen-model)
[![Go Report Card](https://goreportcard.com/badge/github.com/DaoYoung/gen-model)](https://goreportcard.com/report/github.com/DaoYoung/gen-model)
[![GoDoc](https://godoc.org/github.com/xxjwxc/gormt?status.svg)](https://godoc.org/github.com/xxjwxc/gormt)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) 

[中文教程](https://www.jianshu.com/p/0d1d942d281e)
# gen-model
### features
* generate model struct from mysql tables
* make struct fields mappers, and keep at local yaml or mysql table
* support iteration. when table changes, you can update model by reconnect table. when output require fields changes, you can update model by change mapper

### install
```
go get -u github.com/DaoYoung/gen-model
```
### usage
1. run init commend, and set config value at {your_project_dir}/.gen-model.yaml
```
cd ${your_project_dir}
gen-model init
```
2. generate model struct from mysql tables
```
gen-model create

# add model name suffix with `--modelSuffix` or `-m`
# you can set gen.modelSuffix at .gen-model.yaml, but it will covered when set args after commend
gen-model create -m=VO
```
3. create local mappers for struct 
```
gen-model create -y=local-mapper
# it will fail, when run after step 2, because struct file is already exist, it's avoid to cover whole file. you can set `-f=true` to cover it.

gen-model create -y=local-mapper -f=true
```
4. delete one field value in yaml mapper, then recreate struct by it
```
gen-model create -r=local-mapper -f=true
```
5. now you can persist mappers at local yaml, it also support save at mysql database.
```
gen-model create -y=gen-table -f=true
```
6. see what gen-model can do.
```
gen-model -h
gen-model create -h # commend `create` help
```

---
### have fun, free to use