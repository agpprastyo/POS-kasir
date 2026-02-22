
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { InternalProductsStockHistoryResponse } from "@/lib/api/generated"
import { useStockHistoryQuery } from "@/lib/api/query/products"
import { cn } from "@/lib/utils"
import { format } from "date-fns"
import { id as idLocale } from "date-fns/locale"
import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { Badge } from "@/components/ui/badge"
import { ChevronLeft, ChevronRight } from "lucide-react"
import { useTranslation } from "react-i18next"

interface StockHistoryTableProps {
    productId: string
}

export const StockHistoryTable = ({ productId }: StockHistoryTableProps) => {
    const { t } = useTranslation()
    const [page, setPage] = useState(1)
    const limit = 10

    const { data, isLoading } = useStockHistoryQuery(productId, { page, limit })

    const history = data?.history || []
    const pagination = data?.pagination

    if (isLoading) {
        return <div className="space-y-2">
            <Skeleton className="h-10 w-full" />
            <Skeleton className="h-10 w-full" />
            <Skeleton className="h-10 w-full" />
        </div>
    }

    if (!history || history.length === 0) {
        return <div className="text-center py-8 text-muted-foreground">{t('products.stock_history.no_history', 'Belum ada riwayat stok.')}</div>
    }

    const getChangeTypeBadge = (type: string | undefined) => {
        const map: Record<string, string> = {
            sale: t('products.stock_history.type_sale', 'Penjualan'),
            restock: t('products.stock_history.type_restock', 'Stok Masuk'),
            correction: t('products.stock_history.type_correction', 'Koreksi'),
            return: t('products.stock_history.type_return', 'Retur'),
            damage: t('products.stock_history.type_damage', 'Rusak/Hilang')
        }

        let variant: "default" | "secondary" | "destructive" | "outline" = "outline"
        if (type === "sale") variant = "default"
        if (type === "restock") variant = "secondary"
        if (type === "return") variant = "secondary"
        if (type === "correction") variant = "outline"
        if (type === "damage") variant = "destructive"

        return <Badge variant={variant}>{map[type || ""] || type}</Badge>
    }

    return (
        <div className="space-y-4">
            <div className="rounded-md border">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>{t('products.stock_history.date', 'Waktu')}</TableHead>
                            <TableHead>{t('products.stock_history.type', 'Tipe')}</TableHead>
                            <TableHead>{t('products.stock_history.amount', 'Jumlah')}</TableHead>
                            <TableHead>{t('products.stock_history.initial_stock', 'Stok Awal')}</TableHead>
                            <TableHead>{t('products.stock_history.final_stock', 'Stok Akhir')}</TableHead>
                            <TableHead>{t('products.stock_history.note', 'Catatan')}</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {history.map((item: InternalProductsStockHistoryResponse) => (
                            <TableRow key={item.id}>
                                <TableCell>
                                    {item.created_at ? format(new Date(item.created_at), "dd MMM yyyy HH:mm", { locale: idLocale }) : "-"}
                                </TableCell>
                                <TableCell>
                                    {getChangeTypeBadge(item.change_type)}
                                </TableCell>
                                <TableCell>
                                    <span className={cn(
                                        (item.change_amount || 0) > 0 ? "text-primary" : "text-destructive"
                                    )}>
                                        {(item.change_amount || 0) > 0 ? `+${item.change_amount}` : item.change_amount}
                                    </span>
                                </TableCell>
                                <TableCell>{item.previous_stock}</TableCell>
                                <TableCell className="font-bold">{item.current_stock}</TableCell>
                                <TableCell className="max-w-[200px] truncate" title={item.note || ""}>{item.note || "-"}</TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </div>

            <div className="flex items-center justify-end space-x-2 py-4">
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage((old) => Math.max(old - 1, 1))}
                    disabled={page === 1}
                >
                    <ChevronLeft className="h-4 w-4 mr-2" />
                    Previous
                </Button>
                <div className="text-sm">
                    Page {pagination?.current_page || 1} of {pagination?.total_page || 1}
                </div>
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage((old) => old + 1)}
                    disabled={page >= (pagination?.total_page || 1)}
                >
                    Next
                    <ChevronRight className="h-4 w-4 ml-2" />
                </Button>
            </div>
        </div>
    )
}
