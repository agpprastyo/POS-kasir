import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useShiftContext } from "@/context/ShiftContext";
import { useEndShift } from "@/hooks/useShift";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { POSKasirInternalDtoShiftResponse } from "@/lib/api/generated";

const endShiftSchema = z.object({
    actualCashEnd: z.coerce.number().min(0, "Actual cash must be non-negative"),
    password: z.string().min(1, "Password is required"),
});

type EndShiftForm = z.infer<typeof endShiftSchema>;

export const CloseShiftModal: React.FC = () => {
    const { closeShiftModalOpen, setCloseShiftModalOpen } = useShiftContext();
    const { mutate: endShift, isPending } = useEndShift();
    const [summary, setSummary] = useState<POSKasirInternalDtoShiftResponse | null>(null);

    const { register, handleSubmit, formState: { errors }, reset } = useForm<EndShiftForm>({
        resolver: zodResolver(endShiftSchema) as any,
        defaultValues: {
            actualCashEnd: 0,
            password: ""
        }
    });

    const onSubmit = (data: EndShiftForm) => {
        endShift({ actual_cash_end: data.actualCashEnd, password: data.password }, {
            onSuccess: (data) => {
                setSummary(data);
                reset();
            }
        });
    };

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
                            <span className={`font-bold ${isShort ? "text-red-500" : isOver ? "text-green-500" : ""}`}>
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
                <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                    <div className="grid w-full gap-1.5">
                        <Label htmlFor="actualCashEnd">Actual Cash</Label>
                        <Input
                            id="actualCashEnd"
                            type="number"
                            {...register("actualCashEnd")}
                            placeholder="0"
                        />
                        {errors.actualCashEnd && (
                            <p className="text-sm text-red-500">{errors.actualCashEnd.message}</p>
                        )}
                    </div>
                    <div className="grid w-full gap-1.5">
                        <Label htmlFor="password">Password</Label>
                        <Input
                            id="password"
                            type="password"
                            {...register("password")}
                            placeholder=""
                        />
                        {errors.password && (
                            <p className="text-sm text-red-500">{errors.password.message}</p>
                        )}
                    </div>
                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={() => setCloseShiftModalOpen(false)}>Cancel</Button>
                        <Button type="submit" disabled={isPending}>
                            {isPending ? "Closing..." : "Close Register"}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
};
