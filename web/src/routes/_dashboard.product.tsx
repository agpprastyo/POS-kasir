import { createFileRoute } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table'
import { PlusCircle } from 'lucide-react'
import {ReactNode} from "react";
import {FileRouteByToPath} from "@tanstack/router-core/src/routeInfo.ts";

export const Route = createFileRoute('/_dashboard/product' as FileRouteByToPath<any, any>)({
    component: ProductPage,
})

function ProductPage() {
    const products = [
        { id: 'INV001', name: 'Macbook Pro 14"', category: 'Laptop', price: '$1,999', status: 'In Stock' },
        { id: 'INV002', name: 'Iphone 15 Pro', category: 'Smartphone', price: '$999', status: 'Out of Stock' },
        { id: 'INV003', name: 'Dell XPS 13', category: 'Laptop', price: '$1,299', status: 'In Stock' },
        { id: 'INV004', name: 'Samsung Galaxy S24', category: 'Smartphone', price: '$899', status: 'In Stock' },
    ]

    return (
        <div className="flex flex-col gap-4">
            <div className="flex items-center justify-between">
                <h1 className="text-lg font-semibold md:text-2xl">Products</h1>
                <Button size="sm" className="h-8 gap-1">
                    <PlusCircle className="h-3.5 w-3.5" />
                    <span className="sr-only sm:not-sr-only sm:whitespace-nowrap">
            Add Product
          </span>
                </Button>
            </div>

            <div className="rounded-md border">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>Name</TableHead>
                            <TableHead>Category</TableHead>
                            <TableHead>Status</TableHead>
                            <TableHead className="text-right">Price</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {products.map((product) => (
                            <TableRow key={product.id}>
                                <TableCell className="font-medium">{product.name}</TableCell>
                                <TableCell>{product.category}</TableCell>
                                <TableCell>{product.status}</TableCell>
                                <TableCell className="text-right">{product.price}</TableCell>
                            </TableRow>
                        )) as ReactNode}
                    </TableBody>
                </Table>
            </div>
        </div>
    )
}