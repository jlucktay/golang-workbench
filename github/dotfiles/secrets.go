package main

import "go.jlucktay.dev/golang-workbench/secrets"

var (
	ghpaToken = secrets.ReadTokenFromSecrets("./secrets.json")
)
