import React from "react";
import { useForm } from "@tanstack/react-form";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useShiftContext } from "@/context/ShiftContext";
import { useStartShift } from "@/hooks/useShift";
import * as z from "zod";

const startShiftSchema = z.object({
    startCash: z.number().min(0, "Start cash must be non-negative"),
    password: z.string().min(1, "Password is required"),
});


export const OpenShiftModal: React.FC = () => {
    const { openShiftModalOpen, setOpenShiftModalOpen } = useShiftContext();
    const { mutate: startShift, isPending } = useStartShift();

    const form = useForm({
        defaultValues: {
            startCash: 0,
            password: ""
        },
        validators: {
            onChange: startShiftSchema
        },
        onSubmit: async ({ value }) => {
            startShift({ start_cash: value.startCash, password: value.password }, {
                onSuccess: () => {
                    setOpenShiftModalOpen(false);
                    form.reset();
                }
            });
        }
    });

    return (
        <Dialog open={openShiftModalOpen} onOpenChange={setOpenShiftModalOpen}>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Open Register</DialogTitle>
                    <DialogDescription>
                        Enter the starting cash amount and your password to begin.
                    </DialogDescription>
                </DialogHeader>
                <form onSubmit={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    form.handleSubmit();
                }} className="space-y-4">
                    <form.Field
                        name="startCash"
                        children={(field) => (
                            <div className="grid w-full gap-1.5">
                                <Label htmlFor={field.name}>Start Cash</Label>
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
                        <Button type="button" variant="outline" onClick={() => setOpenShiftModalOpen(false)}>Cancel</Button>
                        <form.Subscribe
                            selector={(state) => [state.canSubmit, state.isSubmitting]}
                            children={([canSubmit, isSubmitting]) => (
                                <Button type="submit" disabled={!canSubmit || isSubmitting || isPending}>
                                    {isPending ? "Opening..." : "Open Register"}
                                </Button>
                            )}
                        />
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
};
