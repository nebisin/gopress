package main

import (
	"github.com/nebisin/gopress/controllers"
)

var handler = controllers.Handler{}

func  main()  {
	handler.Initialize()

	handler.Run(":8080")
}