import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { usePrinterSettingsQuery, useUpdatePrinterSettingsMutation, useTestPrintMutation } from "@/lib/api/query/settings"
import { InternalSettingsUpdatePrinterSettingsRequestPaperWidthEnum, InternalSettingsUpdatePrinterSettingsRequestPrintMethodEnum } from "@/lib/api/generated"
import { Loader2, Printer, Bluetooth, Network } from "lucide-react"
import { useReducer, useEffect } from "react"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"

import { printerService } from "@/lib/printer"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

interface PrinterSettingsState {
    connection: string;
    paperWidth: string;
    autoPrint: boolean;
    printMethod: string;
    feConnected: boolean;
    feDeviceName?: string;
    feInterface: "bluetooth" | "lan";
    feLanAddress: string;
}

type PrinterSettingsAction =
    | { type: 'SET_FIELD', field: keyof PrinterSettingsState, value: any }
    | { type: 'MERGE_STATE', payload: Partial<PrinterSettingsState> }

function reducer(state: PrinterSettingsState, action: PrinterSettingsAction): PrinterSettingsState {
    switch (action.type) {
        case 'SET_FIELD':
            return { ...state, [action.field]: action.value }
        case 'MERGE_STATE':
            return { ...state, ...action.payload }
        default:
            return state
    }
}

export function PrinterSettingsCard() {
    const { t } = useTranslation()
    const { data: settings, isLoading } = usePrinterSettingsQuery()

    if (isLoading) {
        return (
            <Card>
                <CardContent className="flex justify-center p-8">
                    <Loader2 className="h-6 w-6 animate-spin" />
                </CardContent>
            </Card>
        )
    }

    return <PrinterSettingsForm settings={settings} t={t} />
}

function PrinterSettingsForm({ settings, t }: { settings: any, t: any }) {
    const updateMutation = useUpdatePrinterSettingsMutation()
    const testPrintMutation = useTestPrintMutation()

    const conn = settings?.connection || ""
    const initialPrintMethod = settings?.print_method || "BE"
    let initialFeInterface: "bluetooth" | "lan" = "bluetooth"
    let initialFeLanAddress = ""
    if (initialPrintMethod === "FE" && conn.startsWith("lan://")) {
        initialFeInterface = "lan"
        initialFeLanAddress = conn.replace("lan://", "")
    }

    const [state, dispatch] = useReducer(reducer, {
        connection: initialPrintMethod === "FE" ? "" : conn,
        paperWidth: settings?.paper_width || "58mm",
        autoPrint: settings?.auto_print || false,
        printMethod: initialPrintMethod,
        feConnected: printerService.isConnected(),
        feDeviceName: printerService.getDeviceName(),
        feInterface: initialFeInterface,
        feLanAddress: initialFeLanAddress
    })

    useEffect(() => {
        const checkConnection = () => {
            const isConn = printerService.isConnected()
            const devName = printerService.getDeviceName()
            if (state.feConnected !== isConn || state.feDeviceName !== devName) {
                dispatch({ type: 'MERGE_STATE', payload: { feConnected: isConn, feDeviceName: devName } })
            }
        }

        checkConnection()
        const interval = setInterval(checkConnection, 2000)
        return () => clearInterval(interval)
    }, [state.feConnected, state.feDeviceName])

    const handleConnectPrinter = async () => {
        try {
            await printerService.connect()
            dispatch({ type: 'MERGE_STATE', payload: { feConnected: true, feDeviceName: printerService.getDeviceName() } })
            toast.success(t('settings.printer.connected', { defaultValue: 'Printer connected' }))
        } catch (error) {
            console.error(error)
            toast.error(t('settings.printer.connection_failed', { defaultValue: 'Failed to connect printer' }))
        }
    }

    const handleSave = () => {
        let finalConnection = state.connection

        if (state.printMethod === "FE") {
            if (state.feInterface === "lan") {
                finalConnection = `lan://${state.feLanAddress}`
            } else {
                finalConnection = "bt://"
            }
        }

        updateMutation.mutate({
            connection: finalConnection,
            paper_width: state.paperWidth as InternalSettingsUpdatePrinterSettingsRequestPaperWidthEnum,
            auto_print: state.autoPrint,
            // @ts-ignore
            print_method: state.printMethod as InternalSettingsUpdatePrinterSettingsRequestPrintMethodEnum
        }, {
            onSuccess: () => {
                toast.success(t('settings.printer.update_success'))
            },
            onError: () => {
                toast.error(t('settings.printer.update_error'))
            }
        })
    }

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

                {state.printMethod === "BE" && (
                    <div className="space-y-2">
                        <Label htmlFor="connection">{t('settings.printer.connection')}</Label>
                        <Input
                            id="connection"
                            value={state.connection}
                            onChange={(e) => dispatch({ type: 'SET_FIELD', field: 'connection', value: e.target.value })}
                            placeholder={t('settings.printer.connection_placeholder')}
                        />
                        <p className="text-xs text-muted-foreground">{t('settings.printer.help_text')}</p>
                    </div>
                )}

                {state.printMethod === "FE" && (
                    <div className="space-y-4 border p-4 rounded-lg bg-muted/20">
                        <Label>{t('settings.printer.interface_type', { defaultValue: 'Interface Type' })}</Label>

                        <Tabs value={state.feInterface} onValueChange={(v) => dispatch({ type: 'SET_FIELD', field: 'feInterface', value: v as any })} className="w-full">
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
                                        <div className={`h-3 w-3 rounded-full ${state.feConnected ? 'bg-primary' : 'bg-destructive'}`} />
                                        <span className="font-medium text-sm">
                                            {state.feConnected
                                                ? `${t('settings.printer.connected_to', { defaultValue: 'Connected to' })}: ${state.feDeviceName || 'Unknown Device'}`
                                                : t('settings.printer.not_connected', { defaultValue: 'Not Connected' })
                                            }
                                        </span>
                                    </div>
                                    <Button size="sm" variant={state.feConnected ? "outline" : "default"} onClick={handleConnectPrinter}>
                                        <Bluetooth className="h-4 w-4 mr-2" />
                                        {state.feConnected ? t('settings.printer.reconnect', { defaultValue: 'Reconnect' }) : t('settings.printer.connect', { defaultValue: 'Connect Printer' })}
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
                                    value={state.feLanAddress}
                                    onChange={(e) => dispatch({ type: 'SET_FIELD', field: 'feLanAddress', value: e.target.value })}
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
                        <Select value={state.paperWidth} onValueChange={(v) => dispatch({ type: 'SET_FIELD', field: 'paperWidth', value: v })}>
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
                            checked={state.autoPrint}
                            onCheckedChange={(v) => dispatch({ type: 'SET_FIELD', field: 'autoPrint', value: v })}
                        />
                    </div>
                </div>

                <div className="space-y-2">
                    <Label>{t('settings.printer.print_method', { defaultValue: 'Print Method' })}</Label>
                    <Select value={state.printMethod} onValueChange={(v) => dispatch({ type: 'SET_FIELD', field: 'printMethod', value: v })}>
                        <SelectTrigger>
                            <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="BE">{t('settings.printer.method_be', { defaultValue: 'Backend (Server connects to Printer)' })}</SelectItem>
                            <SelectItem value="FE">{t('settings.printer.method_fe', { defaultValue: 'Frontend (Browser connects to Printer)' })}</SelectItem>
                        </SelectContent>
                    </Select>
                    <p className="text-xs text-muted-foreground">
                        {state.printMethod === "BE"
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
