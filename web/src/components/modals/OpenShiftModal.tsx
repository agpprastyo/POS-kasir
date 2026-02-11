import React from "react";
import { useForm } from "react-hook-form";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useShiftContext } from "@/context/ShiftContext";
import { useStartShift } from "@/hooks/useShift";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";

const startShiftSchema = z.object({
    startCash: z.coerce.number().min(0, "Start cash must be non-negative"),
    password: z.string().min(1, "Password is required"),
});

type StartShiftForm = z.infer<typeof startShiftSchema>;

export const OpenShiftModal: React.FC = () => {
    const { openShiftModalOpen, setOpenShiftModalOpen } = useShiftContext();
    const { mutate: startShift, isPending } = useStartShift();

    const { register, handleSubmit, formState: { errors }, reset } = useForm<StartShiftForm>({
        resolver: zodResolver(startShiftSchema) as any,
        defaultValues: {
            startCash: 0,
            password: ""
        }
    });

    const onSubmit = (data: StartShiftForm) => {
        startShift({ start_cash: data.startCash, password: data.password }, {
            onSuccess: () => {
                setOpenShiftModalOpen(false);
                reset();
            }
        });
    };

    return (
        <Dialog open={openShiftModalOpen} onOpenChange={setOpenShiftModalOpen}>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Open Register</DialogTitle>
                    <DialogDescription>
                        Enter the starting cash amount and your password to begin.
                    </DialogDescription>
                </DialogHeader>
                <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                    <div className="grid w-full gap-1.5">
                        <Label htmlFor="startCash">Start Cash</Label>
                        <Input
                            id="startCash"
                            type="number"
                            {...register("startCash")}
                            placeholder="0"
                        />
                        {errors.startCash && (
                            <p className="text-sm text-red-500">{errors.startCash.message}</p>
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
                        <Button type="button" variant="outline" onClick={() => setOpenShiftModalOpen(false)}>Cancel</Button>
                        <Button type="submit" disabled={isPending}>
                            {isPending ? "Opening..." : "Open Register"}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
};
