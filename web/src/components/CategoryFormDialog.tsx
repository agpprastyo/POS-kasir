import { InternalCategoriesCreateCategoryRequest } from "@/lib/api/generated"
import { Category, useCreateCategoryMutation, useUpdateCategoryMutation } from "@/lib/api/query/categories"
import { Dialog, DialogContent, DialogTitle, DialogDescription } from "@radix-ui/react-dialog"
import { Loader2 } from "lucide-react"
import { useForm } from '@tanstack/react-form'
import * as z from 'zod'
import { useEffect } from "react"
import { useTranslation } from "react-i18next"

import { DialogHeader, DialogFooter } from "./ui/dialog"
import { Input } from "./ui/input"
import { Button } from "./ui/button"
import { Label } from "./ui/label"


export function CategoryFormDialog({ open, onOpenChange, categoryToEdit }: {
    open: boolean,
    onOpenChange: (open: boolean) => void,
    categoryToEdit: Category | null
}) {
    const { t } = useTranslation()
    const createMutation = useCreateCategoryMutation()
    const updateMutation = useUpdateCategoryMutation()

    const form = useForm({
        defaultValues: {
            name: ''
        },
        validators: {
            onChange: z.object({
                name: z.string().min(1)
            })
        },
        onSubmit: async ({ value }) => {
            const payload: InternalCategoriesCreateCategoryRequest = {
                name: value.name
            }

            try {
                if (categoryToEdit && categoryToEdit.id) {
                    await updateMutation.mutateAsync({ id: categoryToEdit.id, body: payload })
                } else {
                    await createMutation.mutateAsync(payload)
                }
                onOpenChange(false)
                form.reset()
            } catch (error) {
                console.error(error)
            }
        }
    })

    useEffect(() => {
        if (open) {
            if (categoryToEdit) {
                form.setFieldValue('name', categoryToEdit.name || '')
            } else {
                form.reset()
            }
        }
    }, [open, categoryToEdit])

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>{categoryToEdit ? t('settings.category.form.title_edit') : t('settings.category.form.title_add')}</DialogTitle>
                    <DialogDescription>
                        {categoryToEdit ? t('settings.category.form.desc_edit') : t('settings.category.form.desc_add')}
                    </DialogDescription>
                </DialogHeader>
                <form onSubmit={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    form.handleSubmit();
                }}>
                    <div className="grid gap-4 py-4">
                        <form.Field
                            name="name"
                            children={(field) => (
                                <div className="grid gap-2">
                                    <Label htmlFor={field.name}>{t('settings.category.form.name_label')}</Label>
                                    <Input
                                        id={field.name}
                                        name={field.name}
                                        value={field.state.value}
                                        onBlur={field.handleBlur}
                                        onChange={(e) => field.handleChange(e.target.value)}
                                        placeholder={t('settings.category.form.name_placeholder')}
                                    />
                                    {field.state.meta.errors.length > 0 && (
                                        <em role="alert" className="text-[0.8rem] font-medium text-destructive">
                                            {field.state.meta.errors.join(', ')}
                                        </em>
                                    )}
                                </div>
                            )}
                        />
                    </div>
                    <DialogFooter>
                        <form.Subscribe
                            selector={(state) => [state.canSubmit, state.isSubmitting]}
                            children={([canSubmit, isSubmitting]) => (
                                <Button type="submit" disabled={!canSubmit || isSubmitting || createMutation.isPending || updateMutation.isPending}>
                                    {isSubmitting || createMutation.isPending || updateMutation.isPending ? (
                                        <><Loader2 className="mr-2 h-4 w-4 animate-spin" /> {t('settings.category.form.saving')}</>
                                    ) : (
                                        t('settings.category.form.save')
                                    )}
                                </Button>
                            )}
                        />
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}