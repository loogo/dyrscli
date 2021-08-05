# dyrscli

## Description

东易日盛命令行工具

## Prerequire

SystemARCH: amd64

Golang 1.13+

## Get Started

### build library

```
$ .\build.bat
```

### execute library

1. restart Debezium Server tasks
```
$ dyrscli task restart -target=all -host=172.16.10.247:8083
```