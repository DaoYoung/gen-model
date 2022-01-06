[![Go Report Card](https://goreportcard.com/badge/github.com/DaoYoung/gen-model)](https://goreportcard.com/report/github.com/DaoYoung/gen-model)
[![codecov](https://codecov.io/gh/DaoYoung/gen-model/branch/master/graph/badge.svg)](https://codecov.io/gh/DaoYoung/gen-model)
<!--[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) -->

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
1. run `init` command, you will see `.gen-model.yaml`
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
![wHZ6SO.md.gif](https://s1.ax1x.com/2020/09/21/wHZ6SO.md.gif)

4. rename mapper file from `${struct}FieldMapper.yaml` to `${struct}VOFieldMapper.yaml`, and delete one line after fields
```
gen-model create --source=local-mapper --forceCover=true --modelSuffix=VO
# it will generate `${struct}VO.go`
```

![2.gif](https://i.loli.net/2020/09/21/tomFTWGSUyKZNra.gif)

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