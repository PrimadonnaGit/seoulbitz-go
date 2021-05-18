package main

import (
	"fmt"

	router "github.com/PrimadonnaGit/seoulbitz-go/router"
)

func main() {
	// debug mode on and off
	debug := true

	// router를 정의한 파일을 Import 후, router 선언
	echoR := router.Router()

	fmt.Println("Start echo server..")

	if debug {
		// 일반적인 http 서버 실행
		echoR.Logger.Fatal(echoR.Start(":80"))
	} else {
		// 보안 접속을 위한 https 서버 실행
		// "cert.pem"과 "privkey.pem" 파일이 필요함
		echoR.Logger.Fatal(echoR.StartTLS(":80", "cert.pem", "privkey.pem"))
	}
}
