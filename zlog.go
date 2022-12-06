package zlog

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

//type sugaredLogger struct {
//	*zap.SugaredLogger
//	init bool
//}

var sugaredLogger *zap.SugaredLogger

//func init() {
//sugaredLogger = NewLogger("log.log", "debug", true, 10, 10, 10)
//}

func NewLogger(fileName, level string, mod bool, maxSize, maxBackups, maxAge int) {
	//if sugaredLogger.init {
	//	sugaredLogger.Error("has init sugaredLogger")
	//	return sugaredLogger
	//}
	zapLevel, ok := levelMap[level]
	if !ok { //默认info等级输出
		zapLevel = zapcore.InfoLevel
		level = "Hello"
	}
	//fileName := filepath.Join(logDir, logFile)
	_, err := os.Stat(fileName)
	if err != nil {
		file, err := os.Create(fileName)
		if err == nil {
			file.Close()
		} else {
			fmt.Println("Cannot create logFile:", fileName)
			os.Exit(1)
		}

	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//encoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	//	enc.AppendString(cases.Title(language.English).String(l.String())) //自定义level文字输出格式
	//}
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig) //返回文本
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   true,       // 是否压缩
	})

	if zapLevel == zapcore.DebugLevel { //开发者模式，多路输出
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writeSyncer)
	}
	core := zapcore.NewCore(encoder, writeSyncer, zapLevel)

	var log *zap.Logger
	if mod { //调试时，输出具体行号
		log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	} else {
		log = zap.New(core)
	}
	sugaredLogger = log.Sugar()
	//sugaredLogger.init = true
	//return sugaredLogger
}

func Json(arg interface{}) {
	body, _ := json.Marshal(arg)
	sugaredLogger.Debug(string(body))
}
func Debug(args ...interface{}) {
	sugaredLogger.Debug(args...)
}
func Debugf(template string, args ...interface{}) {
	sugaredLogger.Debugf(template, args...)
}
func Info(args ...interface{}) {
	sugaredLogger.Info(args...)
}
func Infof(template string, args ...interface{}) {
	sugaredLogger.Infof(template, args...)
}
func Warn(args ...interface{}) {
	sugaredLogger.Warn(args...)
}
func Warnf(template string, args ...interface{}) {
	sugaredLogger.Warnf(template, args...)
}
func Error(args ...interface{}) {
	sugaredLogger.Error(args...)
}
func Errorf(template string, args ...interface{}) {
	sugaredLogger.Errorf(template, args...)
}
func DPanic(args ...interface{}) {
	sugaredLogger.DPanic(args...)
}
func DPanicf(template string, args ...interface{}) {
	sugaredLogger.DPanicf(template, args...)
}
func Panic(args ...interface{}) {
	sugaredLogger.Panic(args...)
}
func Panicf(template string, args ...interface{}) {
	sugaredLogger.Panicf(template, args...)
}
func Fatal(args ...interface{}) {
	sugaredLogger.Fatal(args...)
}
func Fatalf(template string, args ...interface{}) {
	sugaredLogger.Fatalf(template, args...)
}
