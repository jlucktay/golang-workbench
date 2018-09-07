package main

import "github.com/jlucktay/golang-workbench/secrets"

var (
	ghpaToken = secrets.ReadTokenFromSecrets("./secrets.json")
)
