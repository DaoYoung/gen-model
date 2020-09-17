[![Build Status](https://travis-ci.com/DaoYoung/gen-model.svg?branch=master)](https://travis-ci.com/DaoYoung/gen-model)
[![Go Report Card](https://goreportcard.com/badge/github.com/DaoYoung/gen-model)](https://goreportcard.com/report/github.com/DaoYoung/gen-model)
[![codecov](https://codecov.io/gh/DaoYoung/gen-model/branch/master/graph/badge.svg)](https://codecov.io/gh/DaoYoung/gen-model)
<!--[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) -->

[中文教程](https://www.jianshu.com/p/0d1d942d281e)
# gen-model
### Features
* generate struct file by DB tables
* one table has many struct
* persist mappers from table column to struct attributes
### Install
```
go get -u github.com/DaoYoung/gen-model
```
### Usage
1. run `init` commend, you will see `.gen-model.yaml`
```
cd ${your_project_dir}
gen-model init
```
2. change `mysql.*` `gen.searchTableName` value in `.gen-model.yaml`, then run
```
gen-model create
```
3. create local mappers for struct
```
gen-model create --persist=local-mapper

# it will fail, when run after step 2, because struct file is already exist, it's avoid to cover whole file. you can set `-f=true` to cover it.

gen-model create --persist=local-mapper -f=true
```
![image](https://upload.cc/i1/2020/09/17/76h3Ue.gif)

4. rename mapper file from `${struct}FieldMapper.yaml` to `${struct}VOFieldMapper.yaml`, and delete one line after fields
```
gen-model create --source=local-mapper --forceCover=true --modelSuffix=VO
# it will generate `${struct}VO.go`
```

![image](https://upload.cc/i1/2020/09/17/K6YzdM.gif)

5. persist mapper data in database just use `db-mapper` instead of `local-mapper`.
```
gen-model create --persist=db-mapper --forceCover=true
# it means, you can manage struct for multiple project.
# this require mysql `Create` privilege
```
6. see what gen-model can do.
```
gen-model -h
gen-model create -h # commend `create` help
```