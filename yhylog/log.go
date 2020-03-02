package yhylog

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func LogInit(fileName string){
	logFileName := flag.String("log", fileName, fileName)

	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)


	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "cServer start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//write log
	log.Printf("Rpc Server log ready!!!!!!!!!!!  %s\n", fmt.Sprintf("time=%s", time.Now().String()))
}