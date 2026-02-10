import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useBrandingSettingsQuery, useUpdateBrandingSettingsMutation, useUpdateLogoMutation } from "@/lib/api/query/settings"
import { Loader2, Upload, Trash2 } from "lucide-react"
import { useState, useEffect } from "react"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"

export function BrandingSettingsCard() {
    const { t } = useTranslation()
    const { data: settings, isLoading } = useBrandingSettingsQuery()
    const updateMutation = useUpdateBrandingSettingsMutation()
    const uploadLogoMutation = useUpdateLogoMutation()

    const [appName, setAppName] = useState("")
    const [footerText, setFooterText] = useState("")
    const [themeColor, setThemeColor] = useState("")
    const [themeColorDark, setThemeColorDark] = useState("")
    const [logoUrl, setLogoUrl] = useState("")

    useEffect(() => {
        if (settings) {
            setAppName(settings.app_name || "")
            setFooterText(settings.footer_text || "")
            setThemeColor(settings.theme_color || "")
            setThemeColorDark(settings.theme_color_dark || "")
            setLogoUrl(settings.app_logo || "")
        }
    }, [settings])

    const handleSave = () => {
        updateMutation.mutate({
            app_name: appName,
            footer_text: footerText,
            theme_color: themeColor,
            theme_color_dark: themeColorDark,
            app_logo: logoUrl
        }, {
            onSuccess: () => {
                toast.success(t('settings.branding.update_success'))
            },
            onError: () => {
                toast.error(t('settings.branding.update_error'))
            }
        })
    }

    const handleLogoUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0]
        if (!file) return

        const promise = uploadLogoMutation.mutateAsync(file)
            .then((res) => {
                setLogoUrl(res.url)

            })

        toast.promise(promise, {
            loading: t('settings.branding.uploading'),
            success: t('settings.branding.upload_success'),
            error: t('settings.branding.upload_error')
        })
    }

    const handleRemoveLogo = () => {
        setLogoUrl("")

    }

    if (isLoading) {
        return <div className="flex justify-center p-8"><Loader2 className="h-6 w-6 animate-spin" /></div>
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle>{t('settings.branding.title')}</CardTitle>
                <CardDescription>{t('settings.branding.description')}</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
                <div className="grid w-full items-center gap-1.5">
                    <Label htmlFor="appName">{t('settings.branding.app_name')}</Label>
                    <Input
                        id="appName"
                        value={appName}
                        onChange={(e) => setAppName(e.target.value)}
                        placeholder="POS Kasir"
                    />
                </div>

                <div className="space-y-2">
                    <Label>{t('settings.branding.logo')}</Label>
                    <div className="flex items-start gap-4">
                        <div className="border rounded-lg p-2 h-24 w-24 flex items-center justify-center bg-muted/50 overflow-hidden relative group">
                            {logoUrl ? (
                                <img src={logoUrl} alt="App Logo" className="max-w-full max-h-full object-contain" />
                            ) : (
                                <span className="text-xs text-muted-foreground">{t('settings.branding.no_logo')}</span>
                            )}
                        </div>
                        <div className="flex-1 space-y-2">
                            <div className="flex gap-2">
                                <Button variant="outline" size="sm" className="relative" disabled={uploadLogoMutation.isPending}>
                                    {uploadLogoMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <Upload className="h-4 w-4 mr-2" />}
                                    {t('settings.branding.upload_button')}
                                    <input
                                        type="file"
                                        className="absolute inset-0 opacity-0 cursor-pointer"
                                        accept="image/*"
                                        onChange={handleLogoUpload}
                                        disabled={uploadLogoMutation.isPending}
                                    />
                                </Button>
                                {logoUrl && (
                                    <Button variant="outline" size="sm" onClick={() => setLogoUrl("")}>
                                        <Trash2 className="h-4 w-4 mr-2" />
                                        {t('settings.branding.remove_button')}
                                    </Button>
                                )}
                            </div>
                            <p className="text-xs text-muted-foreground">
                                {t('settings.branding.logo_help')}
                            </p>
                        </div>
                    </div>
                </div>

                <div className="grid w-full items-center gap-1.5">
                    <Label htmlFor="footerText">{t('settings.branding.footer_text')}</Label>
                    <Input
                        id="footerText"
                        value={footerText}
                        onChange={(e) => setFooterText(e.target.value)}
                        placeholder="© 2024 POS Kasir"
                    />
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="grid w-full items-center gap-1.5">
                        <Label htmlFor="themeColor">{t('settings.branding.theme_color')}</Label>
                        <div className="flex gap-2">
                            <Input
                                id="themeColor"
                                type="color"
                                value={themeColor}
                                onChange={(e) => setThemeColor(e.target.value)}
                                className="w-12 p-1 px-1 h-9"
                            />
                            <Input
                                value={themeColor}
                                onChange={(e) => setThemeColor(e.target.value)}
                                placeholder="#000000"
                                className="flex-1"
                            />
                        </div>
                    </div>

                    <div className="grid w-full items-center gap-1.5">
                        <Label htmlFor="themeColorDark">{t('settings.branding.theme_color_dark')}</Label>
                        <div className="flex gap-2">
                            <Input
                                id="themeColorDark"
                                type="color"
                                value={themeColorDark}
                                onChange={(e) => setThemeColorDark(e.target.value)}
                                className="w-12 p-1 px-1 h-9"
                            />
                            <Input
                                value={themeColorDark}
                                onChange={(e) => setThemeColorDark(e.target.value)}
                                placeholder="#ffffff"
                                className="flex-1"
                            />
                        </div>
                    </div>
                </div>

                <div className="flex justify-end">
                    <Button onClick={handleSave} disabled={updateMutation.isPending}>
                        {updateMutation.isPending && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
                        {t('common.save_changes')}
                    </Button>
                </div>
            </CardContent>
        </Card>
    )
}
