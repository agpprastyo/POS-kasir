import { Product, useDeleteProductMutation } from "@/lib/api/query/products.ts";
import React, { useState } from "react";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger
} from "@/components/ui/dropdown-menu.tsx";
import { Button } from "@/components/ui/button.tsx";
import { List, Loader2, MoreHorizontal, Pencil, Trash2 } from "lucide-react";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle
} from "@/components/ui/alert-dialog.tsx";

export function ProductActions({ product, onEdit }: { product: Product, onEdit?: () => void }) {
    const deleteMutation = useDeleteProductMutation()
    const [showDeleteDialog, setShowDeleteDialog] = useState(false)

    const handleDelete = (e: React.MouseEvent) => {
        e.preventDefault()
        if (product.id) {
            deleteMutation.mutate(product.id, {
                onSuccess: () => setShowDeleteDialog(false)
            })
        }
    }

    return (
        <>
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="h-8 w-8 p-0">
                        <span className="sr-only">Open menu</span>
                        <MoreHorizontal className="h-4 w-4" />
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    <DropdownMenuLabel>Actions</DropdownMenuLabel>
                    <DropdownMenuItem onClick={onEdit}>
                        <Pencil className="mr-2 h-4 w-4" /> Edit Product
                    </DropdownMenuItem>
                    <DropdownMenuItem>
                        <List className="mr-2 h-4 w-4" /> Manage Variants
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                        onSelect={(e) => {
                            e.preventDefault()
                            setShowDeleteDialog(true)
                        }}
                        className="text-red-600 focus:text-red-600 cursor-pointer"
                    >
                        <Trash2 className="mr-2 h-4 w-4" /> Delete Product
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>

            <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Delete Product?</AlertDialogTitle>
                        <AlertDialogDescription>
                            Are you sure you want to delete <span
                                className="font-semibold text-foreground">"{product.name}"</span>?
                            This action cannot be undone.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel disabled={deleteMutation.isPending}>Cancel</AlertDialogCancel>
                        <AlertDialogAction
                            onClick={handleDelete}
                            disabled={deleteMutation.isPending}
                            className="bg-red-600 hover:bg-red-700 focus:ring-red-600"
                        >
                            {deleteMutation.isPending ? (
                                <><Loader2 className="mr-2 h-4 w-4 animate-spin" /> Deleting...</>
                            ) : (
                                "Delete"
                            )}
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    )
}