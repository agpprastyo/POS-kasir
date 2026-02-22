import React, { useState } from "react";
import { useForm } from "@tanstack/react-form";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useShiftContext } from "@/context/ShiftContext";
import { useEndShift } from "@/hooks/useShift";
import * as z from "zod";
import { InternalShiftShiftResponse } from "@/lib/api/generated";

const endShiftSchema = z.object({
    actualCashEnd: z.number().min(0, "Actual cash must be non-negative"),
    password: z.string().min(1, "Password is required"),
});


export const CloseShiftModal: React.FC = () => {
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
                        <DialogTitle>Shift Closed</DialogTitle>
                        <DialogDescription>
                            Shift summary and reconciliation.
                        </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4 py-4">
                        <div className="flex justify-between">
                            <span className="text-muted-foreground">Start Cash</span>
                            <span className="font-medium">{summary.start_cash}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-muted-foreground">Expected Cash</span>
                            <span className="font-medium">{summary.expected_cash_end}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-muted-foreground">Actual Cash</span>
                            <span className="font-medium">{summary.actual_cash_end}</span>
                        </div>
                        <div className="border-t pt-2 flex justify-between">
                            <span className="font-bold">Difference</span>
                            <span className={`font-bold ${isShort ? "text-destructive" : isOver ? "text-primary" : ""}`}>
                                {diff}
                            </span>
                        </div>
                    </div>
                    <DialogFooter>
                        <Button onClick={handleClose}>Close</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        );
    }

    return (
        <Dialog open={closeShiftModalOpen} onOpenChange={setCloseShiftModalOpen}>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Close Register</DialogTitle>
                    <DialogDescription>
                        Count the cash in drawer and enter the amount to close the shift.
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
                                <Label htmlFor={field.name}>Actual Cash</Label>
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
                                    <p className="text-sm text-destructive">{field.state.meta.errors.join(', ')}</p>
                                )}
                            </div>
                        )}
                    />
                    <form.Field
                        name="password"
                        children={(field) => (
                            <div className="grid w-full gap-1.5">
                                <Label htmlFor={field.name}>Password</Label>
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
                                    <p className="text-sm text-destructive">{field.state.meta.errors.join(', ')}</p>
                                )}
                            </div>
                        )}
                    />
                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={() => setCloseShiftModalOpen(false)}>Cancel</Button>
                        <form.Subscribe
                            selector={(state) => [state.canSubmit, state.isSubmitting]}
                            children={([canSubmit, isSubmitting]) => (
                                <Button type="submit" disabled={!canSubmit || isSubmitting || isPending}>
                                    {isPending ? "Closing..." : "Close Register"}
                                </Button>
                            )}
                        />
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
};
