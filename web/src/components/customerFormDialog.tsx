import {
    InternalCustomersCreateCustomerRequest,
    InternalCustomersCustomerResponse,
    InternalCustomersUpdateCustomerRequest,
} from "@/lib/api/generated";
import { useTranslation } from "react-i18next";
import { useCreateCustomerMutation, useUpdateCustomerMutation } from "@/lib/api/query/customers.ts";
import { useEffect } from "react";
import { toast } from "sonner";
import { useForm } from '@tanstack/react-form';
import * as z from 'zod';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle
} from "@/components/ui/dialog.tsx";
import { Label } from "@/components/ui/label.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Loader2 } from "lucide-react";



export function CustomerFormDialog({ open, onOpenChange, customerToEdit }: {
    open: boolean,
    onOpenChange: (open: boolean) => void,
    customerToEdit: InternalCustomersCustomerResponse | null
}) {
    const { t } = useTranslation()
    const createMutation = useCreateCustomerMutation()
    const updateMutation = useUpdateCustomerMutation()

    const form = useForm({
        defaultValues: {
            name: '',
            email: '',
            phone: '',
            address: ''
        },
        validators: {
            onChange: z.object({
                name: z.string().min(1, 'Name is required'),
                email: z.string().email('Invalid email format').or(z.literal('')),
                phone: z.string(),
                address: z.string()
            })
        },
        onSubmit: async ({ value }) => {
            try {
                if (customerToEdit) {
                    const payload: InternalCustomersUpdateCustomerRequest = {
                        name: value.name,
                        email: value.email || undefined,
                        phone: value.phone || undefined,
                        address: value.address || undefined
                    }

                    await updateMutation.mutateAsync({ id: customerToEdit.id!, body: payload })
                } else {
                    const payload: InternalCustomersCreateCustomerRequest = {
                        name: value.name,
                        email: value.email || undefined,
                        phone: value.phone || undefined,
                        address: value.address || undefined
                    }
                    await createMutation.mutateAsync(payload)
                }
                onOpenChange(false)
                form.reset()
            } catch (error: any) {
                console.error(error)
                const msg = error.response?.data?.error || "Unknown error";
                toast.error(t('customers.messages.save_failed', { message: msg }))
            }
        }
    })

    useEffect(() => {
        if (open) {
            if (customerToEdit) {
                form.setFieldValue('name', customerToEdit.name || '')
                form.setFieldValue('email', customerToEdit.email || '')
                form.setFieldValue('phone', customerToEdit.phone || '')
                form.setFieldValue('address', customerToEdit.address || '')
            } else {
                form.reset()
            }
        }
    }, [open, customerToEdit])

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>{customerToEdit ? t('customers.form.title_edit', 'Edit Customer') : t('customers.form.title_create', 'Create Customer')}</DialogTitle>
                    <DialogDescription>
                        {customerToEdit ? t('customers.form.desc_edit', 'Modify customer details') : t('customers.form.desc_create', 'Add a new customer to the pos')}
                    </DialogDescription>
                </DialogHeader>
                <form onSubmit={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    form.handleSubmit();
                }}>
                    <div className="grid gap-4 py-4">
                        {/* Name */}
                        <form.Field
                            name="name"
                            children={(field) => (
                                <div className="grid grid-cols-4 items-center gap-4">
                                    <Label htmlFor={field.name} className="text-right">{t('customers.form.name', 'Name')} *</Label>
                                    <div className="col-span-3 flex flex-col gap-1">
                                        <Input
                                            id={field.name}
                                            name={field.name}
                                            value={field.state.value}
                                            onBlur={field.handleBlur}
                                            onChange={e => field.handleChange(e.target.value)}
                                        />
                                        {field.state.meta.errors.length > 0 && (
                                            <em role="alert" className="text-[0.8rem] font-medium text-destructive">
                                                {field.state.meta.errors.join(', ')}
                                            </em>
                                        )}
                                    </div>
                                </div>
                            )}
                        />
                        {/* Email */}
                        <form.Field
                            name="email"
                            children={(field) => (
                                <div className="grid grid-cols-4 items-center gap-4">
                                    <Label htmlFor={field.name} className="text-right">{t('customers.form.email', 'Email')}</Label>
                                    <div className="col-span-3 flex flex-col gap-1">
                                        <Input
                                            id={field.name}
                                            name={field.name}
                                            type="email"
                                            value={field.state.value}
                                            onBlur={field.handleBlur}
                                            onChange={e => field.handleChange(e.target.value)}
                                        />
                                        {field.state.meta.errors.length > 0 && (
                                            <em role="alert" className="text-[0.8rem] font-medium text-destructive">
                                                {field.state.meta.errors.join(', ')}
                                            </em>
                                        )}
                                    </div>
                                </div>
                            )}
                        />
                        {/* Phone */}
                        <form.Field
                            name="phone"
                            children={(field) => (
                                <div className="grid grid-cols-4 items-center gap-4">
                                    <Label htmlFor={field.name} className="text-right">{t('customers.form.phone', 'Phone')}</Label>
                                    <div className="col-span-3 flex flex-col gap-1">
                                        <Input
                                            id={field.name}
                                            name={field.name}
                                            value={field.state.value}
                                            onBlur={field.handleBlur}
                                            onChange={e => field.handleChange(e.target.value)}
                                        />
                                        {field.state.meta.errors.length > 0 && (
                                            <em role="alert" className="text-[0.8rem] font-medium text-destructive">
                                                {field.state.meta.errors.join(', ')}
                                            </em>
                                        )}
                                    </div>
                                </div>
                            )}
                        />
                        {/* Address */}
                        <form.Field
                            name="address"
                            children={(field) => (
                                <div className="grid grid-cols-4 items-center gap-4">
                                    <Label htmlFor={field.name} className="text-right">{t('customers.form.address', 'Address')}</Label>
                                    <div className="col-span-3 flex flex-col gap-1">
                                        <Input
                                            id={field.name}
                                            name={field.name}
                                            value={field.state.value}
                                            onBlur={field.handleBlur}
                                            onChange={e => field.handleChange(e.target.value)}
                                        />
                                        {field.state.meta.errors.length > 0 && (
                                            <em role="alert" className="text-[0.8rem] font-medium text-destructive">
                                                {field.state.meta.errors.join(', ')}
                                            </em>
                                        )}
                                    </div>
                                </div>
                            )}
                        />

                    </div>
                    <DialogFooter>
                        <form.Subscribe
                            selector={(state) => [state.canSubmit, state.isSubmitting]}
                            children={([canSubmit, isSubmitting]) => (
                                <Button type="submit" disabled={!canSubmit || isSubmitting || createMutation.isPending || updateMutation.isPending}>
                                    {(isSubmitting || createMutation.isPending || updateMutation.isPending) && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                    {customerToEdit ? t('customers.form.save_changes', 'Save changes') : t('customers.form.create_customer', 'Add Customer')}
                                </Button>
                            )}
                        />
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}
