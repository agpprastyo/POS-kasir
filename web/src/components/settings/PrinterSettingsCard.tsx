import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"
import { usePrinterSettingsQuery, useUpdatePrinterSettingsMutation, useTestPrintMutation, useDiscoverPrintersQuery } from "@/lib/api/query/settings"
import { InternalSettingsUpdatePrinterSettingsRequestPaperWidthEnum, InternalSettingsUpdatePrinterSettingsRequestPrintMethodEnum } from "@/lib/api/generated"
import { Loader2, Printer, Bluetooth, Network, RefreshCw, CheckCircle2 } from "lucide-react"
import { useReducer, useEffect, useState } from "react"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"

import { printerService } from "@/lib/printer"

interface PrinterSettingsState {
    connection: string;
    paperWidth: string;
    autoPrint: boolean;
    printMethod: "BE" | "FE";
    feConnected: boolean;
    feDeviceName?: string;
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
    const [isDiscovering, setIsDiscovering] = useState(false)
    const { data: discoveredPrinters, refetch: discoverPrinters } = useDiscoverPrintersQuery(isDiscovering)

    const canEdit = updateMutation.isAllowed

    const conn = settings?.connection || ""
    const initialPrintMethod = settings?.print_method || "BE"

    const [state, dispatch] = useReducer(reducer, {
        connection: conn,
        paperWidth: settings?.paper_width || "58mm",
        autoPrint: settings?.auto_print || false,
        printMethod: initialPrintMethod as "BE" | "FE",
        feConnected: printerService.isConnected(),
        feDeviceName: printerService.getDeviceName(),
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

    const handleDiscover = async () => {
        setIsDiscovering(true)
        try {
            await discoverPrinters()
            toast.success(t('settings.printer.discovery_completed', { defaultValue: 'Printer discovery completed' }))
        } catch (error) {
            toast.error(t('settings.printer.discovery_failed', { defaultValue: 'Failed to discover printers' }))
        } finally {
            setIsDiscovering(false)
        }
    }

    const handleSave = () => {
        updateMutation.mutate({
            connection: state.printMethod === "FE" ? "bt://" : state.connection,
            paper_width: state.paperWidth as InternalSettingsUpdatePrinterSettingsRequestPaperWidthEnum,
            auto_print: state.autoPrint,
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
            <CardHeader className="pb-4">
                <div className="flex justify-between items-center">
                    <div>
                        <CardTitle>{t('settings.printer.title')}</CardTitle>
                        <CardDescription>{t('settings.printer.description')}</CardDescription>
                    </div>
                    <Printer className="h-8 w-8 text-muted-foreground opacity-20" />
                </div>
            </CardHeader>
            <CardContent className="space-y-6">
                {/* Method Selection at the top */}
                <div className="space-y-3">
                    <Label className="text-sm font-bold uppercase tracking-widest text-muted-foreground">
                        {t('settings.printer.print_method', { defaultValue: 'Printer Connection Mode' })}
                    </Label>
                    <Tabs
                        value={state.printMethod}
                        onValueChange={(v) => dispatch({ type: 'SET_FIELD', field: 'printMethod', value: v as "BE" | "FE" })}
                        className="w-full"
                    >
                        <TabsList className="grid w-full grid-cols-2 h-12">
                            <TabsTrigger value="BE" className="flex items-center gap-2">
                                <Network className="h-4 w-4" />
                                <span>{t('settings.printer.method_be', { defaultValue: 'Network (LAN / Wi-Fi)' })}</span>
                            </TabsTrigger>
                            <TabsTrigger value="FE" className="flex items-center gap-2">
                                <Bluetooth className="h-4 w-4" />
                                <span>{t('settings.printer.method_fe', { defaultValue: 'Bluetooth' })}</span>
                            </TabsTrigger>
                        </TabsList>

                        <div className="mt-4">
                            <TabsContent value="BE" className="space-y-4 animate-in fade-in duration-300">
                                <div className="p-4 rounded-lg bg-muted/30 border border-dashed text-center space-y-3">
                                    <div className="flex flex-col items-center gap-2">
                                        <p className="text-sm font-medium">{t('settings.printer.discovery_title', { defaultValue: 'Discover Network Printers' })}</p>
                                        <p className="text-sm text-muted-foreground max-w-xs mx-auto">
                                            {t('settings.printer.discovery_desc', { defaultValue: 'Search for available printers in your local network (Port 9100).' })}
                                        </p>
                                    </div>
                                    <Button
                                        variant="secondary"
                                        size="sm"
                                        onClick={handleDiscover}
                                        disabled={isDiscovering}
                                        className="h-10 px-6"
                                    >
                                        {isDiscovering ? (
                                            <RefreshCw className="h-4 w-4 animate-spin mr-2" />
                                        ) : (
                                            <RefreshCw className="h-4 w-4 mr-2" />
                                        )}
                                        {isDiscovering ? t('common.searching', { defaultValue: 'Searching...' }) : t('settings.printer.scan_now', { defaultValue: 'Scan Local Network' })}
                                    </Button>
                                </div>

                                {discoveredPrinters && discoveredPrinters.length > 0 && (
                                    <div className="space-y-2">
                                        <Label className="text-xs font-bold uppercase text-muted-foreground">{t('settings.printer.discovered_count', { defaultValue: 'Printers Found' })} ({discoveredPrinters.length})</Label>
                                        <div className="grid grid-cols-1 gap-2">
                                            {discoveredPrinters.map((p) => (
                                                <button
                                                    key={p.ip}
                                                    type="button"
                                                    onClick={() => dispatch({ type: 'SET_FIELD', field: 'connection', value: p.ip })}
                                                    className={`flex items-center justify-between p-3 rounded-md border transition-all text-sm ${state.connection === p.ip ? 'border-primary bg-primary/5 ring-1 ring-primary' : 'hover:bg-muted/50'}`}
                                                >
                                                    <div className="flex items-center gap-3">
                                                        <Printer className="h-4 w-4 text-muted-foreground" />
                                                        <div className="text-left font-medium">
                                                            {p.name}
                                                        </div>
                                                    </div>
                                                    {state.connection === p.ip && <CheckCircle2 className="h-4 w-4 text-primary" />}
                                                </button>
                                            ))}
                                        </div>
                                    </div>
                                )}

                                <div className="space-y-2 pt-2">
                                    <Label htmlFor="connection" className="text-sm font-semibold">{t('settings.printer.connection', { defaultValue: 'Manual IP / URI' })}</Label>
                                    <Input
                                        id="connection"
                                        value={state.connection}
                                        onChange={(e) => dispatch({ type: 'SET_FIELD', field: 'connection', value: e.target.value })}
                                        placeholder={t('settings.printer.connection_placeholder')}
                                        className="h-10"
                                        disabled={!canEdit}
                                    />
                                    <p className="text-xs text-muted-foreground">{t('settings.printer.help_text')}</p>
                                </div>
                            </TabsContent>

                            <TabsContent value="FE" className="space-y-4 animate-in fade-in duration-300">
                                <div className="space-y-4 border p-5 rounded-lg bg-muted/20">
                                    <div className="flex items-center justify-between">
                                        <div className="flex items-center gap-2">
                                            <Bluetooth className="h-5 w-5 text-primary" />
                                            <Label className="font-semibold text-sm">{t('settings.printer.bluetooth', { defaultValue: 'Bluetooth Printer' })}</Label>
                                        </div>
                                        <div className="flex items-center gap-2">
                                            <div className={`h-2.5 w-2.5 rounded-full ${state.feConnected ? 'bg-primary animate-pulse' : 'bg-destructive'}`} />
                                            <span className="text-xs uppercase tracking-wider font-bold">
                                                {state.feConnected ? t('common.connected', { defaultValue: 'Connected' }) : t('common.disconnected', { defaultValue: 'Disconnected' })}
                                            </span>
                                        </div>
                                    </div>

                                    <div className="flex items-center justify-between bg-background p-4 rounded-md border shadow-sm">
                                        <span className="text-sm font-medium truncate mr-4">
                                            {state.feConnected ? (state.feDeviceName || 'Unknown Device') : t('settings.printer.not_paired', { defaultValue: 'No printer paired' })}
                                        </span>
                                        <Button size="sm" variant={state.feConnected ? "outline" : "default"} onClick={handleConnectPrinter}>
                                            {state.feConnected ? t('settings.printer.reconnect', { defaultValue: 'Change / Reconnect' }) : t('settings.printer.connect', { defaultValue: 'Connect Printer' })}
                                        </Button>
                                    </div>
                                    <p className="text-xs text-muted-foreground italic">
                                        {t('settings.printer.fe_help', { defaultValue: 'Make sure your Bluetooth printer is on and discoverable by this browser.' })}
                                    </p>
                                </div>
                            </TabsContent>
                        </div>
                    </Tabs>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2">
                    <div className="space-y-2">
                        <Label className="text-sm font-semibold">{t('settings.printer.paper_width')}</Label>
                        <Select disabled={!canEdit} value={state.paperWidth} onValueChange={(v) => dispatch({ type: 'SET_FIELD', field: 'paperWidth', value: v })}>
                            <SelectTrigger className="h-12">
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="58mm">{t('settings.printer.paper_58')}</SelectItem>
                                <SelectItem value="80mm">{t('settings.printer.paper_80')}</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="flex items-center justify-between rounded-lg border p-4 bg-muted/5">
                        <div className="space-y-0.5">
                            <Label className="text-sm font-medium">{t('settings.printer.auto_print')}</Label>
                            <p className="text-xs text-muted-foreground leading-tight">Print receipt automatically after payment</p>
                        </div>
                        <Switch
                            disabled={!canEdit}
                            checked={state.autoPrint}
                            onCheckedChange={(v) => dispatch({ type: 'SET_FIELD', field: 'autoPrint', value: v })}
                        />
                    </div>
                </div>

                <div className="flex flex-col-reverse sm:flex-row justify-between gap-3 pt-4 border-t">
                    <Button variant="ghost" onClick={testPrint} size="sm" className="text-muted-foreground hover:text-foreground">
                        <Printer className="h-4 w-4 mr-2" />
                        {t('settings.printer.test_print')}
                    </Button>
                    {canEdit && (
                        <Button onClick={handleSave} disabled={updateMutation.isPending} className="px-8 h-10 shadow-lg shadow-primary/20">
                            {updateMutation.isPending && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
                            {t('common.save_changes')}
                        </Button>
                    )}
                </div>
            </CardContent>
        </Card >
    )
}
