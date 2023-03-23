package zaplog

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
)
var logger *zap.SugaredLogger

func InitLogConfig(path string)  (*Logger,error) {
	var logger = new(Logger)
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return logger, err
	}
	err = yaml.Unmarshal(yamlFile, logger)
	if err != nil {
		fmt.Println(err.Error())
		return logger, err
	}
	return logger, nil
}

// logpath 日志文件路径
// loglevel 日志级别
func InitLogger(path string) error{
	// 日志分割
	logconf,err := InitLogConfig(path)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("log config:",logconf)

	hook := lumberjack.Logger{
		Filename:   logconf.Filename, // 日志文件路径，默认 os.TempDir()
		MaxSize:    logconf.MaxSize,      // 每个日志文件保存10M，默认 100M
		MaxBackups: logconf.MaxBackups,      // 保留30个备份，默认不限
		MaxAge:     logconf.MaxAge,       // 保留7天，默认不限
		Compress:   logconf.Compress,    // 是否压缩，默认不压缩
	}
	write := zapcore.AddSync(&hook)
	// 设置日志级别
	// debug 可以打印出 info debug warn
	// info  级别可以打印 warn info
	// warn  只能打印 warn
	// debug->info->warn->error
	var level zapcore.Level
	switch logconf.Loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	case "warn":
		level = zap.WarnLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	core := zapcore.NewCore(
		 zapcore.NewConsoleEncoder(encoderConfig),
		//zapcore.NewJSONEncoder(encoderConfig),
		// zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&write)), // 打印到控制台和文件
		write,
		level,
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段,如：添加一个服务器名称
//	filed := zap.Fields(zap.String("serviceName", "serviceName"))
	filed := zap.Fields()
	// 构造日志
	logger = zap.New(core, caller, development, filed).Sugar()
	return nil
}
