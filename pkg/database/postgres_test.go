package database

import (
	"POS-kasir/config"
	"POS-kasir/mocks"
	"testing"
	"testing/fstest"

	"go.uber.org/mock/gomock"
)

func TestNewDatabase_InvalidConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockFieldLogger(ctrl)

	cfg := &config.AppConfig{}

	
	cfg.DB.Host = "invalid-host"
	cfg.DB.Port = "5432"
	cfg.DB.User = "user"
	cfg.DB.Password = "pass"
	cfg.DB.DBName = "dbname"
	cfg.DB.SSLMode = "disable"

	mockFS := fstest.MapFS{}

	_, err := NewDatabase(cfg, mockLogger, mockFS)
	if err == nil {
		t.Error("Expected error due to unreachable host, got nil")
	}
}
