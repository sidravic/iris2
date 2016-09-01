package logger


import (
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	"fmt"
	"github.com/supersid/iris2/constants"
)


var logger *log.Logger

func Init(env string, outputFilePath string) *log.Logger {
	logger = &log.Logger{
		Formatter: &log.TextFormatter{},
		Level: log.DebugLevel,
		Out: SetOutput(env, outputFilePath),
	}

	return logger
}

func GetLogger() *log.Logger {
	return logger
}

func SetOutput(env string, filepath string) io.Writer{
	var outputSource io.Writer
	if env == constants.DEVELOPMENT_ENV || env == constants.TEST_ENV  {
		outputSource =  os.Stderr
	} else {
		filename := fmt.Sprintf("%s.log", env)
		var file io.Writer = createFile(filename)
		outputSource = file
	}

	return outputSource
}

func createFile(filename string) io.Writer{
	fullPath := fmt.Sprintf("./../logs/%s", filename)
	file, err := os.OpenFile(fullPath, os.O_WRONLY | os.O_CREATE, 0755)
	if err != nil {
		log.Debug(err.Error())
		panic(err)
	}
	return file
}


