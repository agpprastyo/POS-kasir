package settings_test

import (
	activitylog_repo "POS-kasir/internal/activitylog/repository"
	"POS-kasir/internal/common"
	"POS-kasir/internal/settings"
	settings_repo "POS-kasir/internal/settings/repository"
	"POS-kasir/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSettingsService_GetBranding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSettingsQuerier(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	service := settings.NewSettingsService(nil, nil, mockRepo, nil, mockLogger)

	ctx := context.Background()

	t.Run("SuccessFull", func(t *testing.T) {
		mockRepo.EXPECT().GetSettings(ctx).Return([]settings_repo.Setting{
			{Key: "app_name", Value: "My POS"},
			{Key: "app_logo", Value: "http://logo"},
			{Key: "footer_text", Value: "Footer"},
			{Key: "unknown", Value: "val"},
		}, nil)

		resp, err := service.GetBranding(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "My POS", resp.AppName)
		assert.Equal(t, "http://logo", resp.AppLogo)
		assert.Equal(t, "Footer", resp.FooterText)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().GetSettings(ctx).Return(nil, errors.New("db error"))
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())

		resp, err := service.GetBranding(ctx)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestSettingsService_UpdateBranding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSettingsQuerier(ctrl)
	mockStore := mocks.NewMockStore(ctrl)
	mockActivity := mocks.NewMockIActivityService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	service := settings.NewSettingsService(mockStore, mockActivity, mockRepo, nil, mockLogger)

	ctx := context.WithValue(context.Background(), "user_id", uuid.New())

	t.Run("SuccessAllFields", func(t *testing.T) {
		req := settings.UpdateBrandingRequest{
			AppName:    "New Name",
			AppLogo:    "New Logo",
			FooterText: "New Footer",
		}

		dbMock, _ := pgxmock.NewConn()
		columns := []string{"key", "value", "description", "updated_at"}
		dbMock.ExpectQuery("INSERT INTO settings").WithArgs("app_name", req.AppName).WillReturnRows(pgxmock.NewRows(columns).AddRow("app_name", req.AppName, nil, time.Now()))
		dbMock.ExpectQuery("INSERT INTO settings").WithArgs("app_logo", req.AppLogo).WillReturnRows(pgxmock.NewRows(columns).AddRow("app_logo", req.AppLogo, nil, time.Now()))
		dbMock.ExpectQuery("INSERT INTO settings").WithArgs("footer_text", req.FooterText).WillReturnRows(pgxmock.NewRows(columns).AddRow("footer_text", req.FooterText, nil, time.Now()))
		dbMock.ExpectCommit()

		mockStore.EXPECT().ExecTx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(pgx.Tx) error) error {
			err := fn(dbMock)
			if err == nil {
				return dbMock.Commit(ctx)
			}
			return err
		})
		
		mockRepo.EXPECT().GetSettings(ctx).Return([]settings_repo.Setting{}, nil)
		mockActivity.EXPECT().Log(ctx, gomock.Any(), activitylog_repo.LogActionTypeUPDATE, activitylog_repo.LogEntityTypeSETTINGS, "settings", gomock.Any())

		_, err := service.UpdateBranding(ctx, req)
		assert.NoError(t, err)
		assert.NoError(t, dbMock.ExpectationsWereMet())
	})

	t.Run("PartialFields", func(t *testing.T) {
		req := settings.UpdateBrandingRequest{AppName: "Only Name"}
		mockStore.EXPECT().ExecTx(ctx, gomock.Any()).Return(nil)
		mockRepo.EXPECT().GetSettings(ctx).Return([]settings_repo.Setting{}, nil)
		mockActivity.EXPECT().Log(ctx, gomock.Any(), activitylog_repo.LogActionTypeUPDATE, activitylog_repo.LogEntityTypeSETTINGS, "settings", gomock.Any())

		_, err := service.UpdateBranding(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("TxError", func(t *testing.T) {
		mockStore.EXPECT().ExecTx(ctx, gomock.Any()).Return(errors.New("tx error"))
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())

		_, err := service.UpdateBranding(ctx, settings.UpdateBrandingRequest{AppName: "Fix"})
		assert.Error(t, err)
	})
}

func TestSettingsService_UpdateLogo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSettingsQuerier(ctrl)
	mockR2 := mocks.NewMockIR2(ctrl)
	mockActivity := mocks.NewMockIActivityService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	service := settings.NewSettingsService(nil, mockActivity, mockRepo, mockR2, mockLogger)

	ctx := context.WithValue(context.Background(), "user_id", uuid.New())
	data := []byte("fake-image")

	t.Run("Success", func(t *testing.T) {
		mockR2.EXPECT().UploadFile(ctx, gomock.Any(), data, "image/png").Return("http://logo.url", nil)
		mockRepo.EXPECT().UpsertSetting(ctx, settings_repo.UpsertSettingParams{
			Key:   "app_logo",
			Value: "http://logo.url",
		}).Return(settings_repo.Setting{}, nil)
		mockActivity.EXPECT().Log(ctx, gomock.Any(), activitylog_repo.LogActionTypeUPDATE, activitylog_repo.LogEntityTypeSETTINGS, "settings", gomock.Any())

		url, err := service.UpdateLogo(ctx, data, "logo.png", "image/png")
		assert.NoError(t, err)
		assert.Equal(t, "http://logo.url", url)
	})

	t.Run("UpsertError", func(t *testing.T) {
		mockR2.EXPECT().UploadFile(ctx, gomock.Any(), data, "image/png").Return("http://url", nil)
		mockRepo.EXPECT().UpsertSetting(ctx, gomock.Any()).Return(settings_repo.Setting{}, errors.New("db error"))
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())

		url, err := service.UpdateLogo(ctx, data, "logo.png", "image/png")
		assert.Error(t, err)
		assert.Empty(t, url)
	})

	t.Run("EmptyData", func(t *testing.T) {
		url, err := service.UpdateLogo(ctx, []byte{}, "logo.png", "image/png")
		assert.ErrorIs(t, err, common.ErrBadRequest)
		assert.Empty(t, url)
	})

	t.Run("R2NotInit", func(t *testing.T) {
		srvNoR2 := settings.NewSettingsService(nil, nil, nil, nil, mockLogger)
		mockLogger.EXPECT().Errorf(gomock.Any())
		url, err := srvNoR2.UpdateLogo(ctx, data, "logo.png", "image/png")
		assert.ErrorIs(t, err, common.ErrInternal)
		assert.Empty(t, url)
	})

	t.Run("UploadError", func(t *testing.T) {
		mockR2.EXPECT().UploadFile(ctx, gomock.Any(), data, "image/png").Return("", errors.New("upload error"))
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())

		url, err := service.UpdateLogo(ctx, data, "logo.png", "image/png")
		assert.Error(t, err)
		assert.Empty(t, url)
	})
}

func TestSettingsService_GetPrinterSettings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSettingsQuerier(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	service := settings.NewSettingsService(nil, nil, mockRepo, nil, mockLogger)

	ctx := context.Background()

	t.Run("SuccessAll", func(t *testing.T) {
		mockRepo.EXPECT().GetSettings(ctx).Return([]settings_repo.Setting{
			{Key: "printer_connection", Value: "tcp://1.1.1.1"},
			{Key: "printer_paper_width", Value: "80mm"},
			{Key: "printer_auto_print", Value: "true"},
			{Key: "printer_method", Value: "network"},
		}, nil)

		resp, err := service.GetPrinterSettings(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "tcp://1.1.1.1", resp.Connection)
		assert.Equal(t, "80mm", resp.PaperWidth)
		assert.True(t, resp.AutoPrint)
		assert.Equal(t, "network", resp.PrintMethod)
	})
	
	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().GetSettings(ctx).Return(nil, errors.New("error"))
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
		_, err := service.GetPrinterSettings(ctx)
		assert.Error(t, err)
	})
}

func TestSettingsService_UpdatePrinterSettings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSettingsQuerier(ctrl)
	mockStore := mocks.NewMockStore(ctrl)
	mockActivity := mocks.NewMockIActivityService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	service := settings.NewSettingsService(mockStore, mockActivity, mockRepo, nil, mockLogger)

	ctx := context.WithValue(context.Background(), "user_id", uuid.New())

	t.Run("SuccessAll", func(t *testing.T) {
		autoPrint := true
		req := settings.UpdatePrinterSettingsRequest{
			Connection: "tcp://2.2.2.2",
			PaperWidth: "80mm",
			AutoPrint:  &autoPrint,
			PrintMethod: "network",
		}

		dbMock, _ := pgxmock.NewConn()
		columns := []string{"key", "value", "description", "updated_at"}
		dbMock.ExpectQuery("INSERT INTO settings").WithArgs("printer_connection", req.Connection).WillReturnRows(pgxmock.NewRows(columns).AddRow("printer_connection", req.Connection, nil, time.Now()))
		dbMock.ExpectQuery("INSERT INTO settings").WithArgs("printer_paper_width", req.PaperWidth).WillReturnRows(pgxmock.NewRows(columns).AddRow("printer_paper_width", req.PaperWidth, nil, time.Now()))
		dbMock.ExpectQuery("INSERT INTO settings").WithArgs("printer_auto_print", "true").WillReturnRows(pgxmock.NewRows(columns).AddRow("printer_auto_print", "true", nil, time.Now()))
		dbMock.ExpectQuery("INSERT INTO settings").WithArgs("printer_method", req.PrintMethod).WillReturnRows(pgxmock.NewRows(columns).AddRow("printer_method", req.PrintMethod, nil, time.Now()))
		dbMock.ExpectCommit()

		mockStore.EXPECT().ExecTx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(pgx.Tx) error) error {
			err := fn(dbMock)
			if err == nil {
				return dbMock.Commit(ctx)
			}
			return err
		})

		mockRepo.EXPECT().GetSettings(ctx).Return([]settings_repo.Setting{}, nil)
		mockActivity.EXPECT().Log(ctx, gomock.Any(), activitylog_repo.LogActionTypeUPDATE, activitylog_repo.LogEntityTypeSETTINGS, "settings", gomock.Any())

		_, err := service.UpdatePrinterSettings(ctx, req)
		assert.NoError(t, err)
		assert.NoError(t, dbMock.ExpectationsWereMet())
	})

	t.Run("PartialAutoPrintFalse", func(t *testing.T) {
		autoPrint := false
		req := settings.UpdatePrinterSettingsRequest{AutoPrint: &autoPrint}
		mockStore.EXPECT().ExecTx(ctx, gomock.Any()).Return(nil)
		mockRepo.EXPECT().GetSettings(ctx).Return([]settings_repo.Setting{}, nil)
		mockActivity.EXPECT().Log(ctx, gomock.Any(), activitylog_repo.LogActionTypeUPDATE, activitylog_repo.LogEntityTypeSETTINGS, "settings", gomock.Any())

		_, err := service.UpdatePrinterSettings(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("TxError", func(t *testing.T) {
		mockStore.EXPECT().ExecTx(ctx, gomock.Any()).Return(errors.New("error"))
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
		_, err := service.UpdatePrinterSettings(ctx, settings.UpdatePrinterSettingsRequest{})
		assert.Error(t, err)
	})
}
