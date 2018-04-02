package util

import (
	"io"
	"log"
	"os"

	"../exception"
)

type (
	Logger struct {
		TRACE   *log.Logger
		INFO    *log.Logger
		WARNING *log.Logger
		ERROR   *log.Logger
		XML     *log.Logger
	}

	FilePath struct {
		Trace   string
		Info    string
		Warning string
		Error   string
		XML     string
	}

	File struct {
		Trace   io.Writer
		Info    io.Writer
		Warning io.Writer
		Error   io.Writer
		XML     io.Writer
	}
)

func InitLoggers(u *Utilities) (*Logger, *exception.ASError) {

	errC := make(chan *exception.ASError)
	loggerC := make(chan *Logger)

	go func() {

		filePaths := getFilePaths(u)
		files, asErr := getFiles(filePaths, u)
		if asErr != nil {
			errC <- asErr
			loggerC <- nil
		} else {
			errC <- nil
			loggers := initLoggers(files)
			loggerC <- loggers
		}

	}()

	asErr := <-errC
	loggers := <-loggerC

	if asErr != nil {
		return nil, asErr
	}

	loggers.XML.Println("<?xml version=\"1.0\"?>")
	return loggers, nil

}

/*
	Private methods
*/

func initLoggers(files *File) *Logger {

	return &Logger{
		TRACE: log.New(files.Trace,
			"TRACE: ",
			log.Ldate|log.Ltime|log.Lshortfile),
		INFO: log.New(files.Info,
			"INFO: ",
			log.Ldate|log.Ltime|log.Lshortfile),
		WARNING: log.New(files.Warning,
			"WARNING: ",
			log.Ldate|log.Ltime|log.Lshortfile),
		ERROR: log.New(files.Error,
			"ERROR: ",
			log.Ldate|log.Ltime|log.Lshortfile),
		XML: log.New(files.XML, "", 0),
	}
}

func getFiles(filePaths *FilePath, u *Utilities) (
	*File, *exception.ASError) {

	traceLogFile, err := os.Create(filePaths.Trace)
	infoLogFile, err1 := os.Create(filePaths.Info)
	warningLogFile, err2 := os.Create(filePaths.Warning)
	errorLogFile, err3 := os.Create(filePaths.Error)
	xmlLogFile, err4 := os.Create(filePaths.XML)

	if err != nil ||
		err1 != nil ||
		err2 != nil ||
		err3 != nil ||
		err4 != nil {

		asError := u.GetError(
			exception.AS00015, "create_log_file_error", nil)
		return nil, asError
	}

	return &File{
		Trace:   traceLogFile,
		Info:    infoLogFile,
		Warning: warningLogFile,
		Error:   errorLogFile,
		XML:     xmlLogFile,
	}, nil
}

func getFilePaths(u *Utilities) *FilePath {

	return &FilePath{
		Trace:   u.GetLogFilePath("trace"),
		Info:    u.GetLogFilePath("info"),
		Warning: u.GetLogFilePath("warning"),
		Error:   u.GetLogFilePath("error"),
		XML:     u.GetLogFilePath("xml"),
	}
}
