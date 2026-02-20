import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { usePrinterSettingsQuery, useUpdatePrinterSettingsMutation, useTestPrintMutation } from "@/lib/api/query/settings"
import { InternalSettingsUpdatePrinterSettingsRequestPaperWidthEnum, InternalSettingsUpdatePrinterSettingsRequestPrintMethodEnum } from "@/lib/api/generated"
import { Loader2, Printer } from "lucide-react"
import { useState, useEffect } from "react"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"

import { printerService } from "@/lib/printer"
import { Bluetooth, Network } from "lucide-react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

export function PrinterSettingsCard() {
    const { t } = useTranslation()
    const { data: settings, isLoading } = usePrinterSettingsQuery()
    const updateMutation = useUpdatePrinterSettingsMutation()

    const [connection, setConnection] = useState("")
    const [paperWidth, setPaperWidth] = useState("")
    const [autoPrint, setAutoPrint] = useState(false)
    const [printMethod, setPrintMethod] = useState("BE") // Default to Backend

    // Frontend Printer State
    const [feConnected, setFeConnected] = useState(false)
    const [feDeviceName, setFeDeviceName] = useState<string | undefined>(undefined)
    const [feInterface, setFeInterface] = useState<"bluetooth" | "lan">("bluetooth")
    const [feLanAddress, setFeLanAddress] = useState("")

    // Check FE connection status periodically or on mount
    useEffect(() => {
        const checkConnection = () => {
            setFeConnected(printerService.isConnected())
            setFeDeviceName(printerService.getDeviceName())
        }

        checkConnection()
        const interval = setInterval(checkConnection, 2000)
        return () => clearInterval(interval)
    }, [])

    const handleConnectPrinter = async () => {
        try {
            await printerService.connect()
            setFeConnected(true)
            setFeDeviceName(printerService.getDeviceName())
            toast.success(t('settings.printer.connected', { defaultValue: 'Printer connected' }))
        } catch (error) {
            console.error(error)
            toast.error(t('settings.printer.connection_failed', { defaultValue: 'Failed to connect printer' }))
        }
    }

    useEffect(() => {
        if (settings) {
            setPaperWidth(settings.paper_width || "58mm")
            setAutoPrint(settings.auto_print || false)
            setPrintMethod(settings.print_method || "BE")

            const conn = settings.connection || ""
            if (settings.print_method === "FE") {
                if (conn.startsWith("lan://")) {
                    setFeInterface("lan")
                    setFeLanAddress(conn.replace("lan://", ""))
                } else {
                    setFeInterface("bluetooth")
                }
            } else {
                setConnection(conn)
            }
        }
    }, [settings])

    const handleSave = () => {
        let finalConnection = connection

        if (printMethod === "FE") {
            if (feInterface === "lan") {
                finalConnection = `lan://${feLanAddress}`
            } else {
                finalConnection = "bt://" // Marker for Bluetooth
            }
        }

        updateMutation.mutate({
            connection: finalConnection,
            paper_width: paperWidth as InternalSettingsUpdatePrinterSettingsRequestPaperWidthEnum,
            auto_print: autoPrint,
            // @ts-ignore
            print_method: printMethod as InternalSettingsUpdatePrinterSettingsRequestPrintMethodEnum
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

                {printMethod === "BE" && (
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
                )}

                {printMethod === "FE" && (
                    <div className="space-y-4 border p-4 rounded-lg bg-muted/20">
                        <Label>{t('settings.printer.interface_type', { defaultValue: 'Interface Type' })}</Label>

                        <Tabs value={feInterface} onValueChange={(v) => setFeInterface(v as any)} className="w-full">
                            <TabsList className="grid w-full grid-cols-2">
                                <TabsTrigger value="bluetooth" className="text-xs">
                                    <Bluetooth className="h-3 w-3 mr-2" /> Bluetooth
                                </TabsTrigger>
                                <TabsTrigger value="lan" className="text-xs">
                                    <Network className="h-3 w-3 mr-2" /> LAN / Network
                                </TabsTrigger>
                            </TabsList>

                            <TabsContent value="bluetooth" className="mt-4 space-y-2">
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center gap-2">
                                        <div className={`h-3 w-3 rounded-full ${feConnected ? 'bg-green-500' : 'bg-red-500'}`} />
                                        <span className="font-medium text-sm">
                                            {feConnected
                                                ? `${t('settings.printer.connected_to', { defaultValue: 'Connected to' })}: ${feDeviceName || 'Unknown Device'}`
                                                : t('settings.printer.not_connected', { defaultValue: 'Not Connected' })
                                            }
                                        </span>
                                    </div>
                                    <Button size="sm" variant={feConnected ? "outline" : "default"} onClick={handleConnectPrinter}>
                                        <Bluetooth className="h-4 w-4 mr-2" />
                                        {feConnected ? t('settings.printer.reconnect', { defaultValue: 'Reconnect' }) : t('settings.printer.connect', { defaultValue: 'Connect Printer' })}
                                    </Button>
                                </div>
                                <p className="text-xs text-muted-foreground">
                                    {t('settings.printer.fe_help', { defaultValue: 'Make sure your Bluetooth printer is on and ready to pair.' })}
                                </p>
                            </TabsContent>

                            <TabsContent value="lan" className="mt-4 space-y-2">
                                <Label htmlFor="fe-lan">{t('settings.printer.printer_ip', { defaultValue: 'Printer IP Address / URL' })}</Label>
                                <Input
                                    id="fe-lan"
                                    value={feLanAddress}
                                    onChange={(e) => setFeLanAddress(e.target.value)}
                                    placeholder="e.g. 192.168.1.200 or http://192.168.1.200:80/print"
                                />
                                <p className="text-[10px] text-muted-foreground">
                                    {t('settings.printer.fe_lan_help', { defaultValue: 'Ensure the printer supports HTTP POST raw printing and allows access from this origin (CORS).' })}
                                </p>
                            </TabsContent>
                        </Tabs>
                    </div>
                )}

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="space-y-2">
                        <Label>{t('settings.printer.paper_width')}</Label>
                        <Select value={paperWidth} onValueChange={setPaperWidth}>
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="58mm">58mm</SelectItem>
                                <SelectItem value="80mm">80mm</SelectItem>
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

                <div className="space-y-2">
                    <Label>{t('settings.printer.print_method', { defaultValue: 'Print Method' })}</Label>
                    <Select value={printMethod} onValueChange={setPrintMethod}>
                        <SelectTrigger>
                            <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="BE">{t('settings.printer.method_be', { defaultValue: 'Backend (Server connects to Printer)' })}</SelectItem>
                            <SelectItem value="FE">{t('settings.printer.method_fe', { defaultValue: 'Frontend (Browser connects to Printer)' })}</SelectItem>
                        </SelectContent>
                    </Select>
                    <p className="text-xs text-muted-foreground">
                        {printMethod === "BE"
                            ? t('settings.printer.method_be_desc', { defaultValue: 'Server sends commands directly to printer via IP/Network.' })
                            : t('settings.printer.method_fe_desc', { defaultValue: 'Browser sends commands via Bluetooth/WebSerial. Requires HTTPS or localhost.' })
                        }
                    </p>
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
        </Card >
    )
}
