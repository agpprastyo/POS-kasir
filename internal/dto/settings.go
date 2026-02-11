package dto

import "time"

type SettingResponse struct {
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description *string   `json:"description,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateSettingRequest struct {
	Value string `json:"value" validate:"required"`
}

type BrandingSettingsResponse struct {
	AppName        string `json:"app_name"`
	AppLogo        string `json:"app_logo"`
	FooterText     string `json:"footer_text"`
	ThemeColor     string `json:"theme_color"`
	ThemeColorDark string `json:"theme_color_dark"`
}

type UpdateBrandingRequest struct {
	AppName        string `json:"app_name"`
	AppLogo        string `json:"app_logo"` // URL or empty
	FooterText     string `json:"footer_text"`
	ThemeColor     string `json:"theme_color"`
	ThemeColorDark string `json:"theme_color_dark"`
}

type PrinterSettingsResponse struct {
	Connection string `json:"connection"`
	PaperWidth string `json:"paper_width"`
	AutoPrint  bool   `json:"auto_print"`
}

type UpdatePrinterSettingsRequest struct {
	Connection string `json:"connection"`
	PaperWidth string `json:"paper_width"`
	AutoPrint  *bool  `json:"auto_print"`
}
