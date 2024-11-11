package main

import (
	"cache-demo/internal/logs"
)

func init() {
	logs.SetLogger()
}

func main() {
	Run()
}