import { useEffect } from 'react'
import { useForm } from '@tanstack/react-form'
import { z } from 'zod'
import { useTranslation } from 'react-i18next'
import { format } from 'date-fns'
import { CalendarIcon, Plus, Trash2 } from 'lucide-react'

import { Button } from '@/components/ui/button'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Calendar } from '@/components/ui/calendar'
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from '@/components/ui/popover'
import { cn } from '@/lib/utils'
import {
    Promotion,
    useCreatePromotionMutation,
    useUpdatePromotionMutation
} from '@/lib/api/query/promotions'
import {
    POSKasirInternalPromotionsRepositoryDiscountType,
    POSKasirInternalPromotionsRepositoryPromotionScope,
    POSKasirInternalPromotionsRepositoryPromotionRuleType,
    POSKasirInternalPromotionsRepositoryPromotionTargetType
} from '@/lib/api/generated'
import { useProductsListQuery } from '@/lib/api/query/products'
import { useCategoriesListQuery } from '@/lib/api/query/categories'

const promotionSchema = z.object({
    name: z.string().min(1, 'Name is required'),
    description: z.string().optional(),
    scope: z.nativeEnum(POSKasirInternalPromotionsRepositoryPromotionScope),
    discount_type: z.nativeEnum(POSKasirInternalPromotionsRepositoryDiscountType),
    discount_value: z.coerce.number().min(0),
    max_discount_amount: z.coerce.number().optional(),
    start_date: z.date(),
    end_date: z.date(),
    is_active: z.boolean().default(true),
    rules: z.array(z.object({
        rule_type: z.nativeEnum(POSKasirInternalPromotionsRepositoryPromotionRuleType),
        rule_value: z.string().min(1, "Value required"),
        description: z.string().optional()
    })).default([]),
    targets: z.array(z.object({
        target_type: z.nativeEnum(POSKasirInternalPromotionsRepositoryPromotionTargetType),
        target_id: z.string().min(1, "ID required")
    })).default([])
})

interface PromotionFormDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    promotionToEdit?: Promotion | null
}

export function PromotionFormDialog({ open, onOpenChange, promotionToEdit }: PromotionFormDialogProps) {
    const { t } = useTranslation()
    const createMutation = useCreatePromotionMutation()
    const updateMutation = useUpdatePromotionMutation()

    const { data: productsData } = useProductsListQuery({ limit: 100 })
    const { data: categoriesData } = useCategoriesListQuery({ limit: 100 })

    const products = productsData?.products || []
    const categories = categoriesData || []

    const form = useForm({
        defaultValues: {
            name: '',
            description: '' as string | undefined,
            scope: POSKasirInternalPromotionsRepositoryPromotionScope.PromotionScopeORDER as POSKasirInternalPromotionsRepositoryPromotionScope,
            discount_type: POSKasirInternalPromotionsRepositoryDiscountType.DiscountTypePercentage as POSKasirInternalPromotionsRepositoryDiscountType,
            discount_value: 0,
            max_discount_amount: 0 as number | undefined,
            start_date: new Date(),
            end_date: new Date(),
            is_active: true,
            rules: [] as {
                rule_type: POSKasirInternalPromotionsRepositoryPromotionRuleType;
                rule_value: string;
                description?: string;
            }[],
            targets: [] as {
                target_type: POSKasirInternalPromotionsRepositoryPromotionTargetType;
                target_id: string;
            }[]
        },
        validators: {
            onChange: promotionSchema as any
        },
        onSubmit: async ({ value }) => {
            const payload = {
                ...value,
                start_date: value.start_date.toISOString(),
                end_date: value.end_date.toISOString(),
            }

            if (promotionToEdit) {
                updateMutation.mutate({ id: promotionToEdit.id, body: payload as any }, {
                    onSuccess: () => onOpenChange(false)
                })
            } else {
                createMutation.mutate(payload as any, {
                    onSuccess: () => onOpenChange(false)
                })
            }
        }
    })

    useEffect(() => {
        if (open) {
            if (promotionToEdit) {
                form.reset({
                    name: promotionToEdit.name,
                    description: promotionToEdit.description || '',
                    scope: promotionToEdit.scope,
                    discount_type: promotionToEdit.discount_type,
                    discount_value: promotionToEdit.discount_value,
                    max_discount_amount: promotionToEdit.max_discount_amount || 0,
                    start_date: new Date(promotionToEdit.start_date),
                    end_date: new Date(promotionToEdit.end_date),
                    is_active: promotionToEdit.is_active,
                    rules: promotionToEdit.rules.map(r => ({
                        rule_type: r.rule_type,
                        rule_value: r.rule_value,
                        description: r.description || ''
                    })),
                    targets: promotionToEdit.targets.map(t => ({
                        target_type: t.target_type,
                        target_id: t.target_id
                    }))
                })
            } else {
                form.reset({
                    name: '',
                    description: '',
                    scope: POSKasirInternalPromotionsRepositoryPromotionScope.PromotionScopeORDER,
                    discount_type: POSKasirInternalPromotionsRepositoryDiscountType.DiscountTypePercentage,
                    discount_value: 0,
                    max_discount_amount: 0,
                    start_date: new Date(),
                    end_date: new Date(),
                    is_active: true,
                    rules: [],
                    targets: []
                })
            }
        }
    }, [open, promotionToEdit, form])

    const isLoading = createMutation.isPending || updateMutation.isPending

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>
                        {promotionToEdit ? t('promotions.form.title_edit') : t('promotions.form.title_create')}
                    </DialogTitle>
                    <DialogDescription>
                        {t('promotions.form.desc')}
                    </DialogDescription>
                </DialogHeader>

                <form onSubmit={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    form.handleSubmit();
                }} className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                        <form.Field
                            name="name"
                            children={(field) => (
                                <div className="space-y-2">
                                    <Label htmlFor={field.name}>{t('promotions.form.name')}</Label>
                                    <Input
                                        id={field.name}
                                        name={field.name}
                                        value={field.state.value}
                                        onBlur={field.handleBlur}
                                        onChange={(e) => field.handleChange(e.target.value)}
                                        placeholder={t('promotions.form.name_placeholder')}
                                    />
                                    {field.state.meta.errors.length > 0 && (
                                        <p className="text-[0.8rem] font-medium text-destructive">{field.state.meta.errors.join(', ')}</p>
                                    )}
                                </div>
                            )}
                        />
                        <form.Field
                            name="scope"
                            children={(field) => (
                                <div className="space-y-2">
                                    <Label htmlFor={field.name}>{t('promotions.form.scope')}</Label>
                                    <Select onValueChange={(val) => field.handleChange(val as POSKasirInternalPromotionsRepositoryPromotionScope)} value={field.state.value}>
                                        <SelectTrigger id={field.name}>
                                            <SelectValue placeholder="Select scope" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value={POSKasirInternalPromotionsRepositoryPromotionScope.PromotionScopeORDER}>{t('promotions.scope.ORDER')}</SelectItem>
                                            <SelectItem value={POSKasirInternalPromotionsRepositoryPromotionScope.PromotionScopeITEM}>{t('promotions.scope.ITEM')}</SelectItem>
                                        </SelectContent>
                                    </Select>
                                    {field.state.meta.errors.length > 0 && (
                                        <p className="text-[0.8rem] font-medium text-destructive">{field.state.meta.errors.join(', ')}</p>
                                    )}
                                </div>
                            )}
                        />
                    </div>

                    <form.Field
                        name="description"
                        children={(field) => (
                            <div className="space-y-2">
                                <Label htmlFor={field.name}>{t('promotions.form.description')}</Label>
                                <Textarea
                                    id={field.name}
                                    name={field.name}
                                    value={field.state.value}
                                    onBlur={field.handleBlur}
                                    onChange={(e) => field.handleChange(e.target.value)}
                                />
                                {field.state.meta.errors.length > 0 && (
                                    <p className="text-[0.8rem] font-medium text-destructive">{field.state.meta.errors.join(', ')}</p>
                                )}
                            </div>
                        )}
                    />

                    <div className="grid grid-cols-2 gap-4">
                        <form.Field
                            name="discount_type"
                            children={(field) => (
                                <div className="space-y-2">
                                    <Label htmlFor={field.name}>{t('promotions.form.discount_type')}</Label>
                                    <Select onValueChange={(val) => field.handleChange(val as POSKasirInternalPromotionsRepositoryDiscountType)} value={field.state.value}>
                                        <SelectTrigger id={field.name}>
                                            <SelectValue />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value={POSKasirInternalPromotionsRepositoryDiscountType.DiscountTypePercentage}>{t('promotions.types.percentage')}</SelectItem>
                                            <SelectItem value={POSKasirInternalPromotionsRepositoryDiscountType.DiscountTypeFixedAmount}>{t('promotions.types.fixed_amount')}</SelectItem>
                                        </SelectContent>
                                    </Select>
                                    {field.state.meta.errors.length > 0 && (
                                        <p className="text-[0.8rem] font-medium text-destructive">{field.state.meta.errors.join(', ')}</p>
                                    )}
                                </div>
                            )}
                        />
                        <form.Field
                            name="discount_value"
                            children={(field) => (
                                <div className="space-y-2">
                                    <Label htmlFor={field.name}>{t('promotions.form.discount_value')}</Label>
                                    <Input
                                        id={field.name}
                                        name={field.name}
                                        type="number"
                                        value={field.state.value}
                                        onBlur={field.handleBlur}
                                        onChange={(e) => field.handleChange(Number(e.target.value))}
                                    />
                                    {field.state.meta.errors.length > 0 && (
                                        <p className="text-[0.8rem] font-medium text-destructive">{field.state.meta.errors.join(', ')}</p>
                                    )}
                                </div>
                            )}
                        />
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                        <form.Field
                            name="start_date"
                            children={(field) => (
                                <div className="flex flex-col space-y-2">
                                    <Label htmlFor={field.name}>{t('promotions.form.start_date')}</Label>
                                    <Popover>
                                        <PopoverTrigger asChild>
                                            <Button
                                                variant={"outline"}
                                                className={cn(
                                                    "w-full pl-3 text-left font-normal",
                                                    !field.state.value && "text-muted-foreground"
                                                )}
                                            >
                                                {field.state.value ? (
                                                    format(field.state.value, "PPP")
                                                ) : (
                                                    <span>Pick a date</span>
                                                )}
                                                <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                                            </Button>
                                        </PopoverTrigger>
                                        <PopoverContent className="w-auto p-0" align="start">
                                            <Calendar
                                                mode="single"
                                                selected={field.state.value}
                                                onSelect={(date) => { if (date) field.handleChange(date) }}
                                                initialFocus
                                            />
                                        </PopoverContent>
                                    </Popover>
                                    {field.state.meta.errors.length > 0 && (
                                        <p className="text-[0.8rem] font-medium text-destructive">{field.state.meta.errors.join(', ')}</p>
                                    )}
                                </div>
                            )}
                        />
                        <form.Field
                            name="end_date"
                            children={(field) => (
                                <div className="flex flex-col space-y-2">
                                    <Label htmlFor={field.name}>{t('promotions.form.end_date')}</Label>
                                    <Popover>
                                        <PopoverTrigger asChild>
                                            <Button
                                                variant={"outline"}
                                                className={cn(
                                                    "w-full pl-3 text-left font-normal",
                                                    !field.state.value && "text-muted-foreground"
                                                )}
                                            >
                                                {field.state.value ? (
                                                    format(field.state.value, "PPP")
                                                ) : (
                                                    <span>Pick a date</span>
                                                )}
                                                <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                                            </Button>
                                        </PopoverTrigger>
                                        <PopoverContent className="w-auto p-0" align="start">
                                            <Calendar
                                                mode="single"
                                                selected={field.state.value}
                                                onSelect={(date) => { if (date) field.handleChange(date) }}
                                                initialFocus
                                            />
                                        </PopoverContent>
                                    </Popover>
                                    {field.state.meta.errors.length > 0 && (
                                        <p className="text-[0.8rem] font-medium text-destructive">{field.state.meta.errors.join(', ')}</p>
                                    )}
                                </div>
                            )}
                        />
                    </div>

                    <form.Field
                        name="is_active"
                        children={(field) => (
                            <div className="flex flex-row items-center justify-between rounded-lg border p-4">
                                <div className="space-y-0.5">
                                    <Label className="text-base">{t('promotions.form.is_active')}</Label>
                                </div>
                                <Switch
                                    checked={field.state.value}
                                    onCheckedChange={field.handleChange}
                                />
                            </div>
                        )}
                    />

                    {/* Rules Section */}
                    <form.Field
                        name="rules"
                        mode="array"
                        children={(field) => (
                            <div className="space-y-4">
                                <div className="flex items-center justify-between">
                                    <h3 className="text-sm font-medium">{t('promotions.form.rules')}</h3>
                                    <Button type="button" variant="outline" size="sm" onClick={() => field.pushValue({
                                        rule_type: POSKasirInternalPromotionsRepositoryPromotionRuleType.PromotionRuleTypeMINIMUMORDERAMOUNT,
                                        rule_value: '0'
                                    })}>
                                        <Plus className="mr-2 h-4 w-4" /> {t('promotions.form.add_rule')}
                                    </Button>
                                </div>
                                {field.state.value.map((_, index) => (
                                    <div key={index} className="flex gap-2 items-end border p-4 rounded-md">
                                        <form.Field
                                            name={`rules[${index}].rule_type`}
                                            children={(typeField) => (
                                                <div className="flex-1 space-y-2">
                                                    <Label className="text-xs">{t('promotions.form.rule_type')}</Label>
                                                    <Select onValueChange={(val) => typeField.handleChange(val as POSKasirInternalPromotionsRepositoryPromotionRuleType)} value={typeField.state.value}>
                                                        <SelectTrigger>
                                                            <SelectValue />
                                                        </SelectTrigger>
                                                        <SelectContent>
                                                            <SelectItem value={POSKasirInternalPromotionsRepositoryPromotionRuleType.PromotionRuleTypeMINIMUMORDERAMOUNT}>Min Order Amount</SelectItem>
                                                            <SelectItem value={POSKasirInternalPromotionsRepositoryPromotionRuleType.PromotionRuleTypeREQUIREDPRODUCT}>Required Product</SelectItem>
                                                            <SelectItem value={POSKasirInternalPromotionsRepositoryPromotionRuleType.PromotionRuleTypeREQUIREDCATEGORY}>Required Category</SelectItem>
                                                        </SelectContent>
                                                    </Select>
                                                </div>
                                            )}
                                        />
                                        <form.Field
                                            name={`rules[${index}].rule_value`}
                                            children={(valField) => {
                                                const ruleType = form.getFieldValue(`rules[${index}].rule_type`) as POSKasirInternalPromotionsRepositoryPromotionRuleType

                                                if (ruleType === POSKasirInternalPromotionsRepositoryPromotionRuleType.PromotionRuleTypeREQUIREDPRODUCT) {
                                                    return (
                                                        <div className="flex-1 space-y-2">
                                                            <Label className="text-xs">{t('promotions.form.rule_value')}</Label>
                                                            <Select onValueChange={valField.handleChange} value={valField.state.value}>
                                                                <SelectTrigger>
                                                                    <SelectValue placeholder="Select Product" />
                                                                </SelectTrigger>
                                                                <SelectContent>
                                                                    {products.map(p => (
                                                                        <SelectItem key={p.id || 'unknown'} value={p.id || ''}>{p.name || 'Unnamed Product'}</SelectItem>
                                                                    ))}
                                                                </SelectContent>
                                                            </Select>
                                                        </div>
                                                    )
                                                }

                                                if (ruleType === POSKasirInternalPromotionsRepositoryPromotionRuleType.PromotionRuleTypeREQUIREDCATEGORY) {
                                                    return (
                                                        <div className="flex-1 space-y-2">
                                                            <Label className="text-xs">{t('promotions.form.rule_value')}</Label>
                                                            <Select onValueChange={valField.handleChange} value={valField.state.value}>
                                                                <SelectTrigger>
                                                                    <SelectValue placeholder="Select Category" />
                                                                </SelectTrigger>
                                                                <SelectContent>
                                                                    {categories.map(c => (
                                                                        <SelectItem key={c.id || 'unknown'} value={String(c.id || '')}>{c.name || 'Unnamed Category'}</SelectItem>
                                                                    ))}
                                                                </SelectContent>
                                                            </Select>
                                                        </div>
                                                    )
                                                }

                                                return (
                                                    <div className="flex-1 space-y-2">
                                                        <Label className="text-xs">{t('promotions.form.rule_value')}</Label>
                                                        <Input
                                                            value={valField.state.value}
                                                            onBlur={valField.handleBlur}
                                                            onChange={(e) => valField.handleChange(e.target.value)}
                                                        />
                                                    </div>
                                                )
                                            }}
                                        />
                                        <Button type="button" variant="ghost" size="icon" onClick={() => field.removeValue(index)}>
                                            <Trash2 className="h-4 w-4 text-destructive" />
                                        </Button>
                                    </div>
                                ))}
                            </div>
                        )}
                    />

                    {/* Targets Section */}
                    <form.Field
                        name="targets"
                        mode="array"
                        children={(field) => (
                            <div className="space-y-4">
                                <div className="flex items-center justify-between">
                                    <h3 className="text-sm font-medium">{t('promotions.form.targets')}</h3>
                                    <Button type="button" variant="outline" size="sm" onClick={() => field.pushValue({
                                        target_type: POSKasirInternalPromotionsRepositoryPromotionTargetType.PromotionTargetTypePRODUCT,
                                        target_id: ''
                                    })}>
                                        <Plus className="mr-2 h-4 w-4" /> {t('promotions.form.add_target')}
                                    </Button>
                                </div>
                                {field.state.value.map((_, index) => (
                                    <div key={index} className="flex gap-2 items-end border p-4 rounded-md">
                                        <form.Field
                                            name={`targets[${index}].target_type`}
                                            children={(typeField) => (
                                                <div className="flex-1 space-y-2">
                                                    <Label className="text-xs">{t('promotions.form.target_type')}</Label>
                                                    <Select onValueChange={(val) => typeField.handleChange(val as POSKasirInternalPromotionsRepositoryPromotionTargetType)} value={typeField.state.value}>
                                                        <SelectTrigger>
                                                            <SelectValue />
                                                        </SelectTrigger>
                                                        <SelectContent>
                                                            <SelectItem value={POSKasirInternalPromotionsRepositoryPromotionTargetType.PromotionTargetTypePRODUCT}>Product</SelectItem>
                                                            <SelectItem value={POSKasirInternalPromotionsRepositoryPromotionTargetType.PromotionTargetTypeCATEGORY}>Category</SelectItem>
                                                        </SelectContent>
                                                    </Select>
                                                </div>
                                            )}
                                        />
                                        <form.Field
                                            name={`targets[${index}].target_id`}
                                            children={(idField) => {
                                                const targetType = form.getFieldValue(`targets[${index}].target_type`) as POSKasirInternalPromotionsRepositoryPromotionTargetType

                                                if (targetType === POSKasirInternalPromotionsRepositoryPromotionTargetType.PromotionTargetTypePRODUCT) {
                                                    return (
                                                        <div className="flex-1 space-y-2">
                                                            <Label className="text-xs">{t('promotions.form.target_id')}</Label>
                                                            <Select onValueChange={idField.handleChange} value={idField.state.value}>
                                                                <SelectTrigger>
                                                                    <SelectValue placeholder="Select Product" />
                                                                </SelectTrigger>
                                                                <SelectContent>
                                                                    {products.map(p => (
                                                                        <SelectItem key={p.id || 'unknown'} value={p.id || ''}>{p.name || 'Unnamed Product'}</SelectItem>
                                                                    ))}
                                                                </SelectContent>
                                                            </Select>
                                                        </div>
                                                    )
                                                }

                                                if (targetType === POSKasirInternalPromotionsRepositoryPromotionTargetType.PromotionTargetTypeCATEGORY) {
                                                    return (
                                                        <div className="flex-1 space-y-2">
                                                            <Label className="text-xs">{t('promotions.form.target_id')}</Label>
                                                            <Select onValueChange={idField.handleChange} value={idField.state.value}>
                                                                <SelectTrigger>
                                                                    <SelectValue placeholder="Select Category" />
                                                                </SelectTrigger>
                                                                <SelectContent>
                                                                    {categories.map(c => (
                                                                        <SelectItem key={c.id || 'unknown'} value={String(c.id || '')}>{c.name || 'Unnamed Category'}</SelectItem>
                                                                    ))}
                                                                </SelectContent>
                                                            </Select>
                                                        </div>
                                                    )
                                                }

                                                return (
                                                    <div className="flex-1 space-y-2">
                                                        <Label className="text-xs">{t('promotions.form.target_id')}</Label>
                                                        <Input
                                                            value={idField.state.value}
                                                            onBlur={idField.handleBlur}
                                                            onChange={(e) => idField.handleChange(e.target.value)}
                                                        />
                                                    </div>
                                                )
                                            }}
                                        />
                                        <Button type="button" variant="ghost" size="icon" onClick={() => field.removeValue(index)}>
                                            <Trash2 className="h-4 w-4 text-destructive" />
                                        </Button>
                                    </div>
                                ))}
                            </div>
                        )}
                    />

                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                            {t('common.cancel')}
                        </Button>
                        <form.Subscribe
                            selector={(state) => [state.canSubmit, state.isSubmitting]}
                            children={([canSubmit, isSubmitting]) => (
                                <Button type="submit" disabled={!canSubmit || isSubmitting || isLoading}>
                                    {isLoading || isSubmitting ? t('promotions.form.saving') : t('promotions.form.save')}
                                </Button>
                            )}
                        />
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}
