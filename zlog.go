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

const (
	_createFileErrMsg = "Cannot create logFile:"
	_hasZLoggerErrMsg = "has init ZLogger."
)

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
}

type ZLogger struct {
	*zap.SugaredLogger
	init bool
}

var zLogger = &ZLogger{}

//func init() {
//	logger, _ := zap.NewProduction()
//	defer logger.Sync() // flushes buffer, if any
//	zLogger.SugaredLogger = logger.Sugar()
//}

func (z *ZLogger) Close() {
	//刷新缓冲区
	zLogger.SugaredLogger.Sync()
}

func NewLogger(fileName, level, env string, maxSize, maxBackups, maxAge int) *ZLogger {
	if zLogger.init {
		zLogger.Error(_hasZLoggerErrMsg)
		return zLogger
	}
	zapLevel, ok := levelMap[level]
	if !ok { //默认info等级输出
		zapLevel = zapcore.InfoLevel
	}
	//fileName := filepath.Join(logDir, logFile)
	_, err := os.Stat(fileName)
	if err != nil {
		file, err := os.Create(fileName)
		if err == nil {
			file.Close()
		} else {
			fmt.Println(_createFileErrMsg, fileName)
			os.Exit(1)
		}
	}
	//encoderConfig := zap.NewProductionEncoderConfig()
	//encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//encoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	//	enc.AppendString(cases.Title(language.English).String(l.String())) //自定义level文字输出格式
	//}
	//encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//	enc.AppendString(t.Format("2006-01-02 15:04:05"))
	//}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	//输出格式
	encoder := zapcore.NewConsoleEncoder(encoderConfig) //plain-text 字符串
	//encoder := zapcore.NewJSONEncoder(encoderConfig)  //json 文本

	var core zapcore.Core
	var log *zap.Logger

	switch env {
	case "dev": //开发者模式
		//控制输出的位置
		writeSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fileName,   // 日志文件路径
			MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: maxBackups, // 日志文件最多保存多少个备份
			MaxAge:     maxAge,     // 文件最多保存多少天
			Compress:   true,       // 是否压缩
		})
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writeSyncer) //多路输出
		core = zapcore.NewCore(encoder, writeSyncer, zapLevel)
		log = zap.New(core)
	case "tong": //自定义模式
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel)
		//writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writeSyncer)//多路输出
		log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) //输出行号
	default: //prod
		//控制输出的位置
		writeSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fileName,   // 日志文件路径
			MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: maxBackups, // 日志文件最多保存多少个备份
			MaxAge:     maxAge,     // 文件最多保存多少天
			Compress:   true,       // 是否压缩
		})
		core = zapcore.NewCore(encoder, writeSyncer, zapLevel)
		log = zap.New(core)

	}

	zLogger.SugaredLogger = log.Sugar()
	zLogger.init = true
	return zLogger
}

func Json(arg interface{}) {
	body, _ := json.Marshal(arg)
	zLogger.SugaredLogger.Debug(string(body))
}
func Debug(args ...interface{}) {
	zLogger.SugaredLogger.Debug(args...)
}
func Debugf(template string, args ...interface{}) {
	zLogger.SugaredLogger.Debugf(template, args...)
}
func Info(args ...interface{}) {
	zLogger.SugaredLogger.Info(args...)
}
func Infof(template string, args ...interface{}) {
	zLogger.SugaredLogger.Infof(template, args...)
}
func Warn(args ...interface{}) {
	zLogger.SugaredLogger.Warn(args...)
}
func Warnf(template string, args ...interface{}) {
	zLogger.SugaredLogger.Warnf(template, args...)
}
func Error(args ...interface{}) {
	zLogger.SugaredLogger.Error(args...)
}
func Errorf(template string, args ...interface{}) {
	zLogger.SugaredLogger.Errorf(template, args...)
}
func DPanic(args ...interface{}) {
	zLogger.SugaredLogger.DPanic(args...)
}
func DPanicf(template string, args ...interface{}) {
	zLogger.SugaredLogger.DPanicf(template, args...)
}
func Panic(args ...interface{}) {
	zLogger.SugaredLogger.Panic(args...)
}
func Panicf(template string, args ...interface{}) {
	zLogger.SugaredLogger.Panicf(template, args...)
}
func Fatal(args ...interface{}) {
	zLogger.SugaredLogger.Fatal(args...)
}
func Fatalf(template string, args ...interface{}) {
	zLogger.SugaredLogger.Fatalf(template, args...)
}
