// --- FORM DIALOG (CREATE / EDIT) ---
import {
    POSKasirInternalDtoCreateUserRequest,
    POSKasirInternalDtoProfileResponse,
    POSKasirInternalDtoUpdateUserRequest,
    UsersGetRoleEnum
} from "@/lib/api/generated";
import {useTranslation} from "react-i18next";
import {useCreateUserMutation, useUpdateUserMutation} from "@/lib/api/query/user.ts";
import React, {useEffect, useState} from "react";
import {toast} from "sonner";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle
} from "@/components/ui/dialog.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Loader2} from "lucide-react";

export function UserFormDialog({open, onOpenChange, userToEdit}: {
    open: boolean,
    onOpenChange: (open: boolean) => void,
    userToEdit: POSKasirInternalDtoProfileResponse | null
}) {
    const {t} = useTranslation()
    const createMutation = useCreateUserMutation()
    const updateMutation = useUpdateUserMutation()

    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: '',
        role: UsersGetRoleEnum.Cashier as UsersGetRoleEnum
    })

    useEffect(() => {
        if (open) {
            if (userToEdit) {
                setFormData({
                    username: userToEdit.username || '',
                    email: userToEdit.email || '',
                    password: '',
                    role: (userToEdit.role as UsersGetRoleEnum) || UsersGetRoleEnum.Cashier
                })
            } else {

                setFormData({
                    username: '',
                    email: '',
                    password: '',
                    role: UsersGetRoleEnum.Cashier
                })
            }
        }
    }, [open, userToEdit])

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        try {
            if (userToEdit) {
                const payload: POSKasirInternalDtoUpdateUserRequest = {
                    username: formData.username,
                    email: formData.email,
                    role: formData.role,

                }

                await updateMutation.mutateAsync({id: userToEdit.id!, body: payload})
                toast.success(t('users.form.success_update'))
            } else {
                // --- CREATE LOGIC ---
                const payload: POSKasirInternalDtoCreateUserRequest = {
                    username: formData.username,
                    email: formData.email,
                    password: formData.password,
                    role: formData.role,
                    is_active: true
                }
                await createMutation.mutateAsync(payload)
                toast.success(t('users.form.success_create'))
            }
            onOpenChange(false)
        } catch (error) {
            console.error(error)
            toast.error(t('users.form.error_save'))
        }
    }

    const isSubmitting = createMutation.isPending || updateMutation.isPending

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>{userToEdit ? t('users.form.title_edit') : t('users.form.title_create')}</DialogTitle>
                    <DialogDescription>
                        {userToEdit ? t('users.form.desc_edit') : t('users.form.desc_create')}
                    </DialogDescription>
                </DialogHeader>
                <form onSubmit={handleSubmit}>
                    <div className="grid gap-4 py-4">
                        {/* Username */}
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="username" className="text-right">{t('users.form.username')}</Label>
                            <Input
                                id="username"
                                value={formData.username}
                                onChange={e => setFormData({...formData, username: e.target.value})}
                                className="col-span-3"
                                required
                            />
                        </div>
                        {/* Email */}
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="email" className="text-right">{t('users.form.email')}</Label>
                            <Input
                                id="email"
                                type="email"
                                value={formData.email}
                                onChange={e => setFormData({...formData, email: e.target.value})}
                                className="col-span-3"
                                required
                            />
                        </div>
                        {/* Role */}
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="role" className="text-right">{t('users.form.role')}</Label>
                            <div className="col-span-3">
                                <Select
                                    value={formData.role}
                                    onValueChange={(val: UsersGetRoleEnum) => setFormData({...formData, role: val})}
                                >
                                    <SelectTrigger>
                                        <SelectValue placeholder={t('users.form.select_role')}/>
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value={UsersGetRoleEnum.Admin}>{t('users.roles.admin')}</SelectItem>
                                        <SelectItem
                                            value={UsersGetRoleEnum.Manager}>{t('users.roles.manager')}</SelectItem>
                                        <SelectItem
                                            value={UsersGetRoleEnum.Cashier}>{t('users.roles.cashier')}</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>
                        </div>

                        {!userToEdit && (
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="password" className="text-right">{t('users.form.password')}</Label>
                                <Input
                                    id="password"
                                    type="password"
                                    placeholder={t('users.form.password_placeholder')}
                                    value={formData.password}
                                    onChange={e => setFormData({...formData, password: e.target.value})}
                                    className="col-span-3"
                                    required
                                    minLength={6}
                                />
                            </div>
                        )}
                    </div>
                    <DialogFooter>
                        <Button type="submit" disabled={isSubmitting}>
                            {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin"/>}
                            {userToEdit ? t('users.form.save_changes') : t('users.form.create_user')}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}