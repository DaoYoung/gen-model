[![Build Status](https://travis-ci.com/DaoYoung/gen-model.svg?branch=master)](https://travis-ci.com/DaoYoung/gen-model)
[![Go Report Card](https://goreportcard.com/badge/github.com/DaoYoung/gen-model)](https://goreportcard.com/report/github.com/DaoYoung/gen-model)
[![codecov](https://codecov.io/gh/DaoYoung/gen-model/branch/master/graph/badge.svg)](https://codecov.io/gh/DaoYoung/gen-model)
<!--[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) -->

[中文教程](https://www.jianshu.com/p/0d1d942d281e)
# gen-model
### Features
* Generate struct file by DB tables
* One table has many struct
* Persist mappers from table column to struct attributes
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
2. change `mysql.*` `gen.searchTableName` value in `.gen-model.yaml`, generate struct from mysql tables
```
gen-model create
```
3. create local mappers for struct 
```
gen-model create --persistType=local-mapper

# it will fail, when run after step 2, because struct file is already exist, it's avoid to cover whole file. you can set `-f=true` to cover it.

gen-model create --persistType=local-mapper --forceCover=true
```
4. rename mapper file from `${your_struct_file_name}FieldMapper.yaml` to `${your_struct_file_name}VOFieldMapper.yaml`, and delete one line after fields in `${your_struct_file_name}VOFieldMapper.yaml`
```
gen-model create --sourceType=local-mapper --forceCover=true --modelSuffix=VO
# it will generate `${your_struct_file_name}VO.go`
```
5. persist mapper data in database just user `db-mapper` instead of `local-mapper`.
```
gen-model create --persistType=db-mapper --forceCover=true
# it means, you can manage struct for multiple project.
# this require mysql `Create` privilege
```
6. see what gen-model can do.
```
gen-model -h
gen-model create -h # commend `create` help
```