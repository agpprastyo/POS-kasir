import { InternalCategoriesCreateCategoryRequest } from "@/lib/api/generated"
import { Category, useCreateCategoryMutation, useUpdateCategoryMutation } from "@/lib/api/query/categories"
import { Dialog, DialogContent, DialogTitle, DialogDescription } from "@radix-ui/react-dialog"
import { Loader2 } from "lucide-react"
import { useState, useEffect } from "react"
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

    const [formData, setFormData] = useState({
        name: ''
    })

    useEffect(() => {
        if (open) {

            if (categoryToEdit) {
                setFormData({
                    name: categoryToEdit.name || ''
                })
            } else {
                setFormData({
                    name: ''
                })
            }
        }
    }, [open, categoryToEdit])

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        const payload: InternalCategoriesCreateCategoryRequest = {
            name: formData.name
        }

        try {
            if (categoryToEdit && categoryToEdit.id) {
                await updateMutation.mutateAsync({ id: categoryToEdit.id, body: payload })
            } else {
                await createMutation.mutateAsync(payload)
            }
            onOpenChange(false)
            setFormData({ name: '' })
        } catch (error) {
            console.error(error)
        }
    }

    const isSubmitting = createMutation.isPending || updateMutation.isPending

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>{categoryToEdit ? t('settings.category.form.title_edit') : t('settings.category.form.title_add')}</DialogTitle>
                    <DialogDescription>
                        {categoryToEdit ? t('settings.category.form.desc_edit') : t('settings.category.form.desc_add')}
                    </DialogDescription>
                </DialogHeader>
                <form onSubmit={handleSubmit}>
                    <div className="grid gap-4 py-4">
                        <div className="grid gap-2">
                            <Label htmlFor="name">{t('settings.category.form.name_label')}</Label>
                            <Input
                                id="name"
                                value={formData.name}
                                onChange={(e) => setFormData({ name: e.target.value })}
                                placeholder={t('settings.category.form.name_placeholder')}
                                required
                            />
                        </div>
                    </div>
                    <DialogFooter>
                        <Button type="submit" disabled={isSubmitting}>
                            {isSubmitting ? (
                                <><Loader2 className="mr-2 h-4 w-4 animate-spin" /> {t('settings.category.form.saving')}</>
                            ) : (
                                t('settings.category.form.save')
                            )}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}