import { printerApi, settingsApi } from "@/lib/api/client";

type PrintMethod = "BE" | "FE";

// Minimal Web Bluetooth types to avoid build errors if @types/web-bluetooth is missing
interface BluetoothDevice extends EventTarget {
    id: string;
    name?: string;
    gatt?: BluetoothRemoteGATTServer;
    addEventListener(type: string, listener: EventListenerOrEventListenerObject, options?: boolean | AddEventListenerOptions): void;
    removeEventListener(type: string, listener: EventListenerOrEventListenerObject, options?: boolean | EventListenerOptions): void;
}

interface BluetoothRemoteGATTServer {
    device: BluetoothDevice;
    connected: boolean;
    connect(): Promise<BluetoothRemoteGATTServer>;
    disconnect(): void;
    getPrimaryService(service: BluetoothServiceUUID): Promise<BluetoothRemoteGATTService>;
}

interface BluetoothRemoteGATTService {
    uuid: string;
    device: BluetoothDevice;
    getCharacteristic(characteristic: BluetoothCharacteristicUUID): Promise<BluetoothRemoteGATTCharacteristic>;
}

interface BluetoothRemoteGATTCharacteristic {
    uuid: string;
    service: BluetoothRemoteGATTService;
    writeValue(value: BufferSource): Promise<void>;
}

type BluetoothServiceUUID = number | string;
type BluetoothCharacteristicUUID = number | string;

interface NavigatorBluetooth {
    bluetooth: {
        requestDevice(options?: RequestDeviceOptions): Promise<BluetoothDevice>;
    }
}

interface RequestDeviceOptions {
    filters?: BluetoothLEScanFilter[];
    optionalServices?: BluetoothServiceUUID[];
    acceptAllDevices?: boolean;
}

interface BluetoothLEScanFilter {
    services?: BluetoothServiceUUID[];
    name?: string;
    namePrefix?: string;
}

class PrinterService {
    private device: BluetoothDevice | null = null;
    private characteristic: BluetoothRemoteGATTCharacteristic | null = null;

    async getSettings(): Promise<{ method: PrintMethod, connection: string }> {
        try {
            const response = await settingsApi.settingsPrinterGet();
            const data = (response.data as any).data;
            return {
                method: (data?.print_method as PrintMethod) || "BE",
                connection: data?.connection || ""
            };
        } catch (error) {
            console.error("Failed to fetch printer settings:", error);
            return { method: "BE", connection: "" };
        }
    }

    async printInvoice(orderId: string): Promise<void> {
        const { method, connection } = await this.getSettings();

        if (method === "BE") {
            await this.printBackend(orderId);
        } else {
            // Frontend Printing
            if (connection.startsWith("lan://")) {
                await this.printFrontendNetwork(orderId, connection.replace("lan://", ""));
            } else {
                await this.printFrontendBluetooth(orderId);
            }
        }
    }

    private async printBackend(orderId: string): Promise<void> {
        await printerApi.ordersIdPrintPost(orderId);
    }

    private async getPrintData(orderId: string): Promise<Uint8Array> {
        const response = await printerApi.ordersIdPrintDataGet(orderId);
        const base64Data = (response.data as any).data?.data;

        if (!base64Data) {
            throw new Error("No print data received");
        }

        const binaryString = window.atob(base64Data as unknown as string);
        const len = binaryString.length;
        const bytes = new Uint8Array(len);
        for (let i = 0; i < len; i++) {
            bytes[i] = binaryString.charCodeAt(i);
        }
        return bytes;
    }

    private async printFrontendNetwork(orderId: string, address: string): Promise<void> {
        try {
            const bytes = await this.getPrintData(orderId);

            // Normalize URL
            let url = address;
            if (!url.startsWith("http://") && !url.startsWith("https://")) {
                url = `http://${url}`;
            }

            // Try to append /print if it's just an IP, as a common guess, or use raw if specified
            // Many thermal printers with web interface might accept POST to root or a specific path.
            // Since we can't do raw TCP, we assume the user provided a valid HTTP endpoint or we try root.

            console.log(`Sending print data to ${url}`);

            // We use a blob to send raw bytes
            const blob = new Blob([bytes as unknown as BlobPart], { type: 'application/octet-stream' });

            await fetch(url, {
                method: 'POST',
                body: blob,
                mode: 'no-cors', // Important: Most printers won't have CORS headers. Opaque response.
                cache: 'no-cache',
            });

            // Since mode is no-cors, we can't actually know if it succeeded, but it won't throw on CORS.
            // We assume success if no network error.

        } catch (error) {
            console.error("Frontend LAN printing failed:", error);
            alert(`Printing to ${address} failed. Check console for details. Ensure browser can access the printer IP.`);
            throw error;
        }
    }

    private async printFrontendBluetooth(orderId: string): Promise<void> {
        try {
            // 1. Get raw data
            const bytes = await this.getPrintData(orderId);

            // 2. Connect if needed
            if (!this.device || !this.device.gatt?.connected) {
                await this.connect();
            }

            // 3. Send data
            if (this.characteristic) {
                // ESC/POS printers usually use chunks to avoid buffer overflow
                const chunkSize = 512;
                for (let i = 0; i < bytes.length; i += chunkSize) {
                    const chunk = bytes.slice(i, i + chunkSize);
                    await this.characteristic.writeValue(chunk);
                    // Small delay between chunks to prevent buffer overflow on printer
                    await new Promise(resolve => setTimeout(resolve, 50));
                }
            } else {
                throw new Error("Printer not connected");
            }

        } catch (error) {
            console.error("Frontend printing failed:", error);
            alert("Printing failed. Please check printer connection."); // Simple user feedback
            throw error;
        }
    }

    async connect(): Promise<void> {
        const nav = navigator as unknown as NavigatorBluetooth;
        if (!nav.bluetooth) {
            throw new Error("Web Bluetooth API not supported in this browser");
        }

        try {
            this.device = await nav.bluetooth.requestDevice({
                filters: [
                    { services: ['000018f0-0000-1000-8000-00805f9b34fb'] } // Standard UUID for 18f0 service usually found in thermal printers
                ],
                optionalServices: ['000018f0-0000-1000-8000-00805f9b34fb']
            });

            if (!this.device.gatt) {
                throw new Error("Device does not support GATT");
            }

            const server = await this.device.gatt.connect();
            const service = await server.getPrimaryService('000018f0-0000-1000-8000-00805f9b34fb');
            this.characteristic = await service.getCharacteristic('00002af1-0000-1000-8000-00805f9b34fb'); // Write characteristic

            this.device.addEventListener('gattserverdisconnected', this.onDisconnected);

        } catch (error) {
            console.error("Bluetooth connection failed:", error);
            throw error;
        }
    }

    private onDisconnected = () => {
        console.log("Printer disconnected");
        this.device = null;
        this.characteristic = null;
    }

    isConnected(): boolean {
        return !!(this.device && this.device.gatt && this.device.gatt.connected);
    }

    getDeviceName(): string | undefined {
        return this.device?.name;
    }
}

export const printerService = new PrinterService();
