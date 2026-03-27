import { formatDate, formatRupiah } from "@/lib/utils"
import { Pencil, Trash2, MoreHorizontal, RotateCcw } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Promotion } from '@/lib/api/query/promotions'
import { POSKasirInternalPromotionsRepositoryDiscountType } from '@/lib/api/generated'

interface PromotionsTableProps {
    promotions: Promotion[]
    t: any
    onEdit: (promo: Promotion) => void
    onDelete: (id: string) => void
    onRestore?: (id: string) => void
    isTrash: boolean
    canEdit?: boolean
    canDelete?: boolean
    canRestore?: boolean
}

export function PromotionsTable({
    promotions,
    t,
    onEdit,
    onDelete,
    onRestore,
    isTrash,
    canEdit,
    canDelete,
    canRestore
}: PromotionsTableProps) {
    const hasAnyAction = (canEdit || canDelete) && !isTrash || (canRestore && isTrash)

    return (
        <div className="rounded-md border bg-card">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>{t('promotions.table.name')}</TableHead>
                        <TableHead>{t('promotions.table.scope')}</TableHead>
                        <TableHead>{t('promotions.table.discount')}</TableHead>
                        <TableHead>{t('promotions.table.period')}</TableHead>
                        <TableHead>{t('promotions.table.status')}</TableHead>
                        {hasAnyAction && <TableHead className="text-right">{t('promotions.table.actions')}</TableHead>}
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {promotions.length === 0 ? (
                        <TableRow>
                            <TableCell colSpan={hasAnyAction ? 6 : 5} className="h-24 text-center text-muted-foreground">
                                {t('promotions.table.empty')}
                            </TableCell>
                        </TableRow>
                    ) : (
                        promotions.map((promo: Promotion) => (
                            <TableRow key={promo.id}>
                                <TableCell className="font-medium">
                                    <div className="flex flex-col">
                                        <span>{promo.name}</span>
                                        {promo.description && <span className="text-sm text-muted-foreground">{promo.description}</span>}
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <Badge variant="outline">{promo.scope}</Badge>
                                </TableCell>
                                <TableCell>
                                    <div className="flex flex-col">
                                        <span className="font-bold">
                                            {promo.discount_type === POSKasirInternalPromotionsRepositoryDiscountType.DiscountTypePercentage
                                                ? `${promo.discount_value}%`
                                                : formatRupiah(promo.discount_value)}
                                        </span>
                                        {promo.max_discount_amount && promo.max_discount_amount > 0 && (
                                            <span className="text-sm text-muted-foreground">{t('promotions.table.max')} {formatRupiah(promo.max_discount_amount)}</span>
                                        )}
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <div className="text-sm">
                                        {formatDate(promo.start_date)} - {formatDate(promo.end_date)}
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <Badge variant={promo.is_active ? 'default' : 'secondary'}>
                                        {promo.is_active ? t('promotions.status.active') : t('promotions.status.inactive')}
                                    </Badge>
                                </TableCell>
                                {hasAnyAction && (
                                    <TableCell className="text-right">
                                        <DropdownMenu>
                                            <DropdownMenuTrigger asChild>
                                                <Button variant="ghost" className="h-8 w-8 p-0">
                                                    <span className="sr-only">{t('promotions.actions.open_menu')}</span>
                                                    <MoreHorizontal className="h-4 w-4" />
                                                </Button>
                                            </DropdownMenuTrigger>
                                            <DropdownMenuContent align="end">
                                                <DropdownMenuLabel>{t('common.actions')}</DropdownMenuLabel>
                                                {!isTrash && (
                                                    <>
                                                        {canEdit && (
                                                            <DropdownMenuItem onClick={() => onEdit(promo)}>
                                                                <Pencil className="mr-2 h-4 w-4" /> {t('common.edit')}
                                                            </DropdownMenuItem>
                                                        )}
                                                        {canEdit && canDelete && <DropdownMenuSeparator />}
                                                        {canDelete && (
                                                            <DropdownMenuItem onClick={() => onDelete(promo.id)} className="text-destructive">
                                                                <Trash2 className="mr-2 h-4 w-4" /> {t('common.delete')}
                                                            </DropdownMenuItem>
                                                        )}
                                                    </>
                                                )}
                                                {isTrash && onRestore && canRestore && (
                                                    <DropdownMenuItem onClick={() => onRestore(promo.id)}>
                                                        <RotateCcw className="mr-2 h-4 w-4" /> {t('common.restore')}
                                                    </DropdownMenuItem>
                                                )}
                                            </DropdownMenuContent>
                                        </DropdownMenu>
                                    </TableCell>
                                )}
                            </TableRow>
                        ))
                    )}
                </TableBody>
            </Table>
        </div>
    )
}
