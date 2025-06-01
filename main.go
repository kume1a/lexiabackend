package main

import (
	"lexia/internal/config"
	_ "time/tzdata"
)

func main() {
	config.LoadEnv()
}
