import React, { useState } from "react";
import { useForm } from "@tanstack/react-form";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useTranslation } from "react-i18next";
import { useShiftContext } from "@/context/ShiftContext";
import { useEndShift } from "@/hooks/useShift";
import * as z from "zod";
import { InternalShiftShiftResponse } from "@/lib/api/generated";

const endShiftSchema = z.object({
    actualCashEnd: z.number().min(0, "Actual cash must be non-negative"),
    password: z.string().min(1, "Password is required"),
});


export const CloseShiftModal: React.FC = () => {
    const { t } = useTranslation();
    const { closeShiftModalOpen, setCloseShiftModalOpen } = useShiftContext();
    const { mutate: endShift, isPending } = useEndShift();
    const [summary, setSummary] = useState<InternalShiftShiftResponse | null>(null);

    const form = useForm({
        defaultValues: {
            actualCashEnd: 0,
            password: ""
        },
        validators: {
            onChange: endShiftSchema
        },
        onSubmit: async ({ value }) => {
            endShift({ actual_cash_end: value.actualCashEnd, password: value.password }, {
                onSuccess: (data) => {
                    setSummary(data);
                    form.reset();
                }
            });
        }
    });

    const handleClose = () => {
        setCloseShiftModalOpen(false);
        setSummary(null);
    };

    if (summary) {
        // Show summary view
        const diff = summary.difference || 0;
        const isShort = diff < 0;
        const isOver = diff > 0;

        return (
            <Dialog open={closeShiftModalOpen} onOpenChange={handleClose}>
                <DialogContent className="sm:max-w-[425px]">
                    <DialogHeader>
                        <DialogTitle>{t('shift.close_modal.title_summary')}</DialogTitle>
                        <DialogDescription>
                            {t('shift.close_modal.desc_summary')}
                        </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4 py-4">
                        <div className="flex justify-between">
                            <span className="text-muted-foreground">{t('shift.close_modal.start_cash')}</span>
                            <span className="font-medium">{summary.start_cash}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-muted-foreground">{t('shift.close_modal.expected_cash')}</span>
                            <span className="font-medium">{summary.expected_cash_end}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-muted-foreground">{t('shift.close_modal.actual_cash')}</span>
                            <span className="font-medium">{summary.actual_cash_end}</span>
                        </div>
                        <div className="border-t pt-2 flex justify-between">
                            <span className="font-bold">{t('shift.close_modal.difference')}</span>
                            <span className={`font-bold ${isShort ? "text-destructive" : isOver ? "text-primary" : ""}`}>
                                {diff}
                            </span>
                        </div>
                    </div>
                    <DialogFooter>
                        <Button onClick={handleClose}>{t('shift.close_modal.close')}</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        );
    }

    return (
        <Dialog open={closeShiftModalOpen} onOpenChange={setCloseShiftModalOpen}>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>{t('shift.close_modal.title')}</DialogTitle>
                    <DialogDescription>
                        {t('shift.close_modal.desc')}
                    </DialogDescription>
                </DialogHeader>
                <form onSubmit={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    form.handleSubmit();
                }} className="space-y-4">
                    <form.Field
                        name="actualCashEnd"
                        children={(field) => (
                            <div className="grid w-full gap-1.5">
                                <Label htmlFor={field.name}>{t('shift.close_modal.actual_cash')}</Label>
                                <Input
                                    id={field.name}
                                    type="number"
                                    name={field.name}
                                    value={field.state.value}
                                    onBlur={field.handleBlur}
                                    onChange={(e) => field.handleChange(Number(e.target.value))}
                                    placeholder="0"
                                />
                                {field.state.meta.errors.length > 0 && (
                                    <p className="text-sm text-destructive">{field.state.meta.errors.map(err => typeof err === 'object' ? ((err as any).message ?? JSON.stringify(err)) : String(err)).join(', ')}</p>
                                )}
                            </div>
                        )}
                    />
                    <form.Field
                        name="password"
                        children={(field) => (
                            <div className="grid w-full gap-1.5">
                                <Label htmlFor={field.name}>{t('shift.close_modal.password')}</Label>
                                <Input
                                    id={field.name}
                                    type="password"
                                    name={field.name}
                                    value={field.state.value}
                                    onBlur={field.handleBlur}
                                    onChange={(e) => field.handleChange(e.target.value)}
                                    placeholder=""
                                />
                                {field.state.meta.errors.length > 0 && (
                                    <p className="text-sm text-destructive">{field.state.meta.errors.map(err => typeof err === 'object' ? ((err as any).message ?? JSON.stringify(err)) : String(err)).join(', ')}</p>
                                )}
                            </div>
                        )}
                    />
                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={() => setCloseShiftModalOpen(false)}>{t('common.cancel')}</Button>
                        <form.Subscribe
                            selector={(state) => [state.canSubmit, state.isSubmitting]}
                            children={([canSubmit, isSubmitting]) => (
                                <Button type="submit" disabled={!canSubmit || isSubmitting || isPending}>
                                    {isPending ? t('shift.close_modal.closing') : t('shift.close_modal.submit')}
                                </Button>
                            )}
                        />
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
};
