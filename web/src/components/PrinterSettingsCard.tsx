import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { usePrinterSettingsQuery, useUpdatePrinterSettingsMutation, useTestPrintMutation } from "@/lib/api/query/settings"
import { Loader2, Printer } from "lucide-react"
import { useState, useEffect } from "react"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"

export function PrinterSettingsCard() {
    const { t } = useTranslation()
    const { data: settings, isLoading } = usePrinterSettingsQuery()
    const updateMutation = useUpdatePrinterSettingsMutation()

    const [connection, setConnection] = useState("")
    const [paperWidth, setPaperWidth] = useState("")
    const [autoPrint, setAutoPrint] = useState(false)

    useEffect(() => {
        if (settings) {
            setConnection(settings.connection || "")
            setPaperWidth(settings.paper_width || "58")
            setAutoPrint(settings.auto_print || false)
        }
    }, [settings])

    const handleSave = () => {
        updateMutation.mutate({
            connection: connection,
            paper_width: paperWidth,
            auto_print: autoPrint
        }, {
            onSuccess: () => {
                toast.success(t('settings.printer.update_success'))
            },
            onError: () => {
                toast.error(t('settings.printer.update_error'))
            }
        })
    }

    const testPrintMutation = useTestPrintMutation()

    const testPrint = () => {
        testPrintMutation.mutate(undefined, {
            onSuccess: () => {
                toast.success(t('settings.printer.test_print_success', { defaultValue: 'Test print sent successfully' }))
            },
            onError: () => {
                toast.error(t('settings.printer.test_print_error', { defaultValue: 'Failed to send test print' }))
            }
        })
    }

    if (isLoading) {
        return <div className="flex justify-center p-8"><Loader2 className="h-6 w-6 animate-spin" /></div>
    }

    return (
        <Card>
            <CardHeader>
                <div className="flex justify-between items-center">
                    <div>
                        <CardTitle>{t('settings.printer.title')}</CardTitle>
                        <CardDescription>{t('settings.printer.description')}</CardDescription>
                    </div>
                    <Printer className="h-8 w-8 text-muted-foreground opacity-20" />
                </div>
            </CardHeader>
            <CardContent className="space-y-6">
                <div className="space-y-2">
                    <Label htmlFor="connection">{t('settings.printer.connection')}</Label>
                    <Input
                        id="connection"
                        value={connection}
                        onChange={(e) => setConnection(e.target.value)}
                        placeholder={t('settings.printer.connection_placeholder')}
                    />
                    <p className="text-xs text-muted-foreground">{t('settings.printer.help_text')}</p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="space-y-2">
                        <Label>{t('settings.printer.paper_width')}</Label>
                        <Select value={paperWidth} onValueChange={setPaperWidth}>
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="58">58mm</SelectItem>
                                <SelectItem value="80">80mm</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="flex items-center justify-between rounded-lg border p-4">
                        <div className="space-y-0.5">
                            <Label className="text-base">{t('settings.printer.auto_print')}</Label>
                        </div>
                        <Switch
                            checked={autoPrint}
                            onCheckedChange={setAutoPrint}
                        />
                    </div>
                </div>

                <div className="flex justify-between pt-4">
                    <Button variant="outline" onClick={testPrint}>
                        {t('settings.printer.test_print')}
                    </Button>
                    <Button onClick={handleSave} disabled={updateMutation.isPending}>
                        {updateMutation.isPending && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
                        {t('common.save_changes')}
                    </Button>
                </div>
            </CardContent>
        </Card>
    )
}
