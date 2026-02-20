package settings

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
	AppName        string `json:"app_name" validate:"required,min=3,max=100"`
	AppLogo        string `json:"app_logo" validate:"omitempty"`
	FooterText     string `json:"footer_text" validate:"omitempty,max=200"`
	ThemeColor     string `json:"theme_color" validate:"required,hexcolor"`
	ThemeColorDark string `json:"theme_color_dark" validate:"required,hexcolor"`
}

type PrinterSettingsResponse struct {
	Connection  string `json:"connection"`
	PaperWidth  string `json:"paper_width"`
	AutoPrint   bool   `json:"auto_print"`
	PrintMethod string `json:"print_method"`
}

type UpdatePrinterSettingsRequest struct {
	Connection  string `json:"connection" validate:"required"`
	PaperWidth  string `json:"paper_width" validate:"required,oneof=58mm 80mm"`
	AutoPrint   *bool  `json:"auto_print" validate:"required"`
	PrintMethod string `json:"print_method" validate:"required,oneof=BE FE"`
}
