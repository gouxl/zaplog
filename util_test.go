package zaplog

import "testing"

func TestInitLogger(t *testing.T) {
	InitLogger("./logger.yaml")
	for i := 0; i < 99; i++ {
		logger.Debug("debug",i)
		logger.Warn("warn")
		logger.Error("error",i)
		logger.Info("info")
	}

}