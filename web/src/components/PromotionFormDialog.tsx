import { useEffect } from 'react'
import { useForm, useFieldArray } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
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
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form'
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
    POSKasirInternalRepositoryDiscountType,
    POSKasirInternalRepositoryPromotionScope,
    POSKasirInternalRepositoryPromotionRuleType,
    POSKasirInternalRepositoryPromotionTargetType
} from '@/lib/api/generated'
import { useProductsListQuery } from '@/lib/api/query/products'
import { useCategoriesListQuery } from '@/lib/api/query/categories'

const promotionSchema = z.object({
    name: z.string().min(1, 'Name is required'),
    description: z.string().optional(),
    scope: z.enum(POSKasirInternalRepositoryPromotionScope),
    discount_type: z.enum(POSKasirInternalRepositoryDiscountType),
    discount_value: z.coerce.number().min(0),
    max_discount_amount: z.coerce.number().optional(),
    start_date: z.date(),
    end_date: z.date(),
    is_active: z.boolean().default(true),
    rules: z.array(z.object({
        rule_type: z.enum(POSKasirInternalRepositoryPromotionRuleType),
        rule_value: z.string().min(1, "Value required"),
        description: z.string().optional()
    })).default([]),
    targets: z.array(z.object({
        target_type: z.enum(POSKasirInternalRepositoryPromotionTargetType),
        target_id: z.string().min(1, "ID required")
    })).default([])
})

type PromotionFormValues = z.infer<typeof promotionSchema>

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

    const form = useForm<PromotionFormValues>({
        resolver: zodResolver(promotionSchema) as any,
        defaultValues: {
            name: '',
            description: '',
            scope: POSKasirInternalRepositoryPromotionScope.PromotionScopeORDER,
            discount_type: POSKasirInternalRepositoryDiscountType.DiscountTypePercentage,
            discount_value: 0,
            max_discount_amount: 0,
            start_date: new Date(),
            end_date: new Date(),
            is_active: true,
            rules: [],
            targets: []
        },
    })

    const { fields: ruleFields, append: appendRule, remove: removeRule } = useFieldArray({
        control: form.control,
        name: "rules"
    })

    const { fields: targetFields, append: appendTarget, remove: removeTarget } = useFieldArray({
        control: form.control,
        name: "targets"
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
                    scope: POSKasirInternalRepositoryPromotionScope.PromotionScopeORDER,
                    discount_type: POSKasirInternalRepositoryDiscountType.DiscountTypePercentage,
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

    const onSubmit = (values: PromotionFormValues) => {
        const payload = {
            ...values,
            start_date: values.start_date.toISOString(),
            end_date: values.end_date.toISOString(),
        }

        if (promotionToEdit) {
            updateMutation.mutate({ id: promotionToEdit.id, body: payload }, {
                onSuccess: () => onOpenChange(false)
            })
        } else {
            createMutation.mutate(payload, {
                onSuccess: () => onOpenChange(false)
            })
        }
    }

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

                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <div className="grid grid-cols-2 gap-4">
                            <FormField
                                control={form.control}
                                name="name"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>{t('promotions.form.name')}</FormLabel>
                                        <FormControl>
                                            <Input placeholder={t('promotions.form.name_placeholder')} {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="scope"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>{t('promotions.form.scope')}</FormLabel>
                                        <Select onValueChange={field.onChange} defaultValue={field.value}>
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Select scope" />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                <SelectItem value={POSKasirInternalRepositoryPromotionScope.PromotionScopeORDER}>{t('promotions.scope.ORDER')}</SelectItem>
                                                <SelectItem value={POSKasirInternalRepositoryPromotionScope.PromotionScopeITEM}>{t('promotions.scope.ITEM')}</SelectItem>
                                            </SelectContent>
                                        </Select>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>

                        <FormField
                            control={form.control}
                            name="description"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>{t('promotions.form.description')}</FormLabel>
                                    <FormControl>
                                        <Textarea {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        <div className="grid grid-cols-2 gap-4">
                            <FormField
                                control={form.control}
                                name="discount_type"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>{t('promotions.form.discount_type')}</FormLabel>
                                        <Select onValueChange={field.onChange} defaultValue={field.value}>
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                <SelectItem value={POSKasirInternalRepositoryDiscountType.DiscountTypePercentage}>{t('promotions.types.percentage')}</SelectItem>
                                                <SelectItem value={POSKasirInternalRepositoryDiscountType.DiscountTypeFixedAmount}>{t('promotions.types.fixed_amount')}</SelectItem>
                                            </SelectContent>
                                        </Select>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="discount_value"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>{t('promotions.form.discount_value')}</FormLabel>
                                        <FormControl>
                                            <Input type="number" {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            <FormField
                                control={form.control}
                                name="start_date"
                                render={({ field }) => (
                                    <FormItem className="flex flex-col">
                                        <FormLabel>{t('promotions.form.start_date')}</FormLabel>
                                        <Popover>
                                            <PopoverTrigger asChild>
                                                <FormControl>
                                                    <Button
                                                        variant={"outline"}
                                                        className={cn(
                                                            "w-full pl-3 text-left font-normal",
                                                            !field.value && "text-muted-foreground"
                                                        )}
                                                    >
                                                        {field.value ? (
                                                            format(field.value, "PPP")
                                                        ) : (
                                                            <span>Pick a date</span>
                                                        )}
                                                        <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                                                    </Button>
                                                </FormControl>
                                            </PopoverTrigger>
                                            <PopoverContent className="w-auto p-0" align="start">
                                                <Calendar
                                                    mode="single"
                                                    selected={field.value}
                                                    onSelect={field.onChange}
                                                    initialFocus
                                                />
                                            </PopoverContent>
                                        </Popover>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="end_date"
                                render={({ field }) => (
                                    <FormItem className="flex flex-col">
                                        <FormLabel>{t('promotions.form.end_date')}</FormLabel>
                                        <Popover>
                                            <PopoverTrigger asChild>
                                                <FormControl>
                                                    <Button
                                                        variant={"outline"}
                                                        className={cn(
                                                            "w-full pl-3 text-left font-normal",
                                                            !field.value && "text-muted-foreground"
                                                        )}
                                                    >
                                                        {field.value ? (
                                                            format(field.value, "PPP")
                                                        ) : (
                                                            <span>Pick a date</span>
                                                        )}
                                                        <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                                                    </Button>
                                                </FormControl>
                                            </PopoverTrigger>
                                            <PopoverContent className="w-auto p-0" align="start">
                                                <Calendar
                                                    mode="single"
                                                    selected={field.value}
                                                    onSelect={field.onChange}
                                                    initialFocus
                                                />
                                            </PopoverContent>
                                        </Popover>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>

                        <FormField
                            control={form.control}
                            name="is_active"
                            render={({ field }) => (
                                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                                    <div className="space-y-0.5">
                                        <FormLabel className="text-base">{t('promotions.form.is_active')}</FormLabel>
                                    </div>
                                    <FormControl>
                                        <Switch
                                            checked={field.value}
                                            onCheckedChange={field.onChange}
                                        />
                                    </FormControl>
                                </FormItem>
                            )}
                        />

                        {/* Rules Section */}
                        <div className="space-y-4">
                            <div className="flex items-center justify-between">
                                <h3 className="text-sm font-medium">{t('promotions.form.rules')}</h3>
                                <Button type="button" variant="outline" size="sm" onClick={() => appendRule({
                                    rule_type: POSKasirInternalRepositoryPromotionRuleType.PromotionRuleTypeMINIMUMORDERAMOUNT,
                                    rule_value: '0'
                                })}>
                                    <Plus className="mr-2 h-4 w-4" /> {t('promotions.form.add_rule')}
                                </Button>
                            </div>
                            {ruleFields.map((field, index) => (
                                <div key={field.id} className="flex gap-2 items-end border p-4 rounded-md">
                                    <FormField
                                        control={form.control}
                                        name={`rules.${index}.rule_type`}
                                        render={({ field }) => (
                                            <FormItem className="flex-1">
                                                <FormLabel className="text-xs">{t('promotions.form.rule_type')}</FormLabel>
                                                <Select onValueChange={field.onChange} defaultValue={field.value}>
                                                    <FormControl>
                                                        <SelectTrigger>
                                                            <SelectValue />
                                                        </SelectTrigger>
                                                    </FormControl>
                                                    <SelectContent>
                                                        <SelectItem value={POSKasirInternalRepositoryPromotionRuleType.PromotionRuleTypeMINIMUMORDERAMOUNT}>Min Order Amount</SelectItem>
                                                        <SelectItem value={POSKasirInternalRepositoryPromotionRuleType.PromotionRuleTypeREQUIREDPRODUCT}>Required Product</SelectItem>
                                                        <SelectItem value={POSKasirInternalRepositoryPromotionRuleType.PromotionRuleTypeREQUIREDCATEGORY}>Required Category</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name={`rules.${index}.rule_value`}
                                        render={({ field }) => {
                                            const ruleType = form.watch(`rules.${index}.rule_type`)

                                            if (ruleType === POSKasirInternalRepositoryPromotionRuleType.PromotionRuleTypeREQUIREDPRODUCT) {
                                                return (
                                                    <FormItem className="flex-1">
                                                        <FormLabel className="text-xs">{t('promotions.form.rule_value')}</FormLabel>
                                                        <Select onValueChange={field.onChange} defaultValue={field.value} value={field.value}>
                                                            <FormControl>
                                                                <SelectTrigger>
                                                                    <SelectValue placeholder="Select Product" />
                                                                </SelectTrigger>
                                                            </FormControl>
                                                            <SelectContent>
                                                                {products.map(p => (
                                                                    <SelectItem key={p.id || 'unknown'} value={p.id || ''}>{p.name || 'Unnamed Product'}</SelectItem>
                                                                ))}
                                                            </SelectContent>
                                                        </Select>
                                                    </FormItem>
                                                )
                                            }

                                            if (ruleType === POSKasirInternalRepositoryPromotionRuleType.PromotionRuleTypeREQUIREDCATEGORY) {
                                                return (
                                                    <FormItem className="flex-1">
                                                        <FormLabel className="text-xs">{t('promotions.form.rule_value')}</FormLabel>
                                                        <Select onValueChange={field.onChange} defaultValue={field.value} value={field.value}>
                                                            <FormControl>
                                                                <SelectTrigger>
                                                                    <SelectValue placeholder="Select Category" />
                                                                </SelectTrigger>
                                                            </FormControl>
                                                            <SelectContent>
                                                                {categories.map(c => (
                                                                    <SelectItem key={c.id || 'unknown'} value={String(c.id || '')}>{c.name || 'Unnamed Category'}</SelectItem>
                                                                ))}
                                                            </SelectContent>
                                                        </Select>
                                                    </FormItem>
                                                )
                                            }

                                            return (
                                                <FormItem className="flex-1">
                                                    <FormLabel className="text-xs">{t('promotions.form.rule_value')}</FormLabel>
                                                    <FormControl>
                                                        <Input {...field} />
                                                    </FormControl>
                                                </FormItem>
                                            )
                                        }}
                                    />
                                    <Button type="button" variant="ghost" size="icon" onClick={() => removeRule(index)}>
                                        <Trash2 className="h-4 w-4 text-destructive" />
                                    </Button>
                                </div>
                            ))}
                        </div>

                        {/* Targets Section (Simplified) */}
                        <div className="space-y-4">
                            <div className="flex items-center justify-between">
                                <h3 className="text-sm font-medium">{t('promotions.form.targets')}</h3>
                                <Button type="button" variant="outline" size="sm" onClick={() => appendTarget({
                                    target_type: POSKasirInternalRepositoryPromotionTargetType.PromotionTargetTypePRODUCT,
                                    target_id: ''
                                })}>
                                    <Plus className="mr-2 h-4 w-4" /> {t('promotions.form.add_target')}
                                </Button>
                            </div>
                            {targetFields.map((field, index) => (
                                <div key={field.id} className="flex gap-2 items-end border p-4 rounded-md">
                                    <FormField
                                        control={form.control}
                                        name={`targets.${index}.target_type`}
                                        render={({ field }) => (
                                            <FormItem className="flex-1">
                                                <FormLabel className="text-xs">{t('promotions.form.target_type')}</FormLabel>
                                                <Select onValueChange={field.onChange} defaultValue={field.value}>
                                                    <FormControl>
                                                        <SelectTrigger>
                                                            <SelectValue />
                                                        </SelectTrigger>
                                                    </FormControl>
                                                    <SelectContent>
                                                        <SelectItem value={POSKasirInternalRepositoryPromotionTargetType.PromotionTargetTypePRODUCT}>Product</SelectItem>
                                                        <SelectItem value={POSKasirInternalRepositoryPromotionTargetType.PromotionTargetTypeCATEGORY}>Category</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name={`targets.${index}.target_id`}
                                        render={({ field }) => {
                                            const targetType = form.watch(`targets.${index}.target_type`)

                                            if (targetType === POSKasirInternalRepositoryPromotionTargetType.PromotionTargetTypePRODUCT) {
                                                return (
                                                    <FormItem className="flex-1">
                                                        <FormLabel className="text-xs">{t('promotions.form.target_id')}</FormLabel>
                                                        <Select onValueChange={field.onChange} defaultValue={field.value} value={field.value}>
                                                            <FormControl>
                                                                <SelectTrigger>
                                                                    <SelectValue placeholder="Select Product" />
                                                                </SelectTrigger>
                                                            </FormControl>
                                                            <SelectContent>
                                                                {products.map(p => (
                                                                    <SelectItem key={p.id || 'unknown'} value={p.id || ''}>{p.name || 'Unnamed Product'}</SelectItem>
                                                                ))}
                                                            </SelectContent>
                                                        </Select>
                                                    </FormItem>
                                                )
                                            }

                                            if (targetType === POSKasirInternalRepositoryPromotionTargetType.PromotionTargetTypeCATEGORY) {
                                                return (
                                                    <FormItem className="flex-1">
                                                        <FormLabel className="text-xs">{t('promotions.form.target_id')}</FormLabel>
                                                        <Select onValueChange={field.onChange} defaultValue={field.value} value={field.value}>
                                                            <FormControl>
                                                                <SelectTrigger>
                                                                    <SelectValue placeholder="Select Category" />
                                                                </SelectTrigger>
                                                            </FormControl>
                                                            <SelectContent>
                                                                {categories.map(c => (
                                                                    <SelectItem key={c.id || 'unknown'} value={String(c.id || '')}>{c.name || 'Unnamed Category'}</SelectItem>
                                                                ))}
                                                            </SelectContent>
                                                        </Select>
                                                    </FormItem>
                                                )
                                            }

                                            return (
                                                <FormItem className="flex-1">
                                                    <FormLabel className="text-xs">{t('promotions.form.target_id')}</FormLabel>
                                                    <FormControl>
                                                        <Input {...field} />
                                                    </FormControl>
                                                </FormItem>
                                            )
                                        }}
                                    />
                                    <Button type="button" variant="ghost" size="icon" onClick={() => removeTarget(index)}>
                                        <Trash2 className="h-4 w-4 text-destructive" />
                                    </Button>
                                </div>
                            ))}
                        </div>

                        <DialogFooter>
                            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                                {t('common.cancel')}
                            </Button>
                            <Button type="submit" disabled={isLoading}>
                                {isLoading ? t('promotions.form.saving') : t('promotions.form.save')}
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    )
}
