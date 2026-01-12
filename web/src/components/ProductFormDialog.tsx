import { useEffect, useRef, useState } from 'react'
import {
    type Product,
    useCreateProductMutation,
    useUpdateProductMutation,
    useUploadProductImageMutation,
    useCreateProductOptionMutation,
    useUpdateProductOptionMutation,
    useUploadProductOptionImageMutation,
    useProductDetailQuery
} from '@/lib/api/query/products'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue, } from "@/components/ui/select"
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { ImageIcon, Loader2, Plus, Trash } from 'lucide-react'
import { POSKasirInternalDtoCreateProductRequest, POSKasirInternalDtoUpdateProductRequest } from "@/lib/api/generated";
import { toast } from "sonner";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle
} from "@/components/ui/dialog.tsx";
import { Label } from "@/components/ui/label.tsx";
import { useImageCropper } from "@/hooks/use-image-cropper.ts";
import { ImageCropperDialog } from "@/components/common/ImageCropperDialog.tsx";


interface VariantItem {
    id?: string
    name: string
    additional_price: number
    image_url?: string | null
    imageFile?: File | null
    isNew?: boolean
}

export function ProductFormDialog({ open, onOpenChange, productToEdit, categories }: {
    open: boolean,
    onOpenChange: (open: boolean) => void,
    productToEdit: Product | null,
    categories: any[]
}) {
    const createMutation = useCreateProductMutation()
    const updateMutation = useUpdateProductMutation()
    const uploadImageMutation = useUploadProductImageMutation()

    const createOptionMutation = useCreateProductOptionMutation()
    const updateOptionMutation = useUpdateProductOptionMutation()
    const uploadOptionImageMutation = useUploadProductOptionImageMutation()

    const { data: detailProduct, isLoading: isLoadingDetail } = useProductDetailQuery(productToEdit?.id || '')

    const cropper = useImageCropper()

    const [formData, setFormData] = useState({
        name: '',
        category_id: 0,
        price: 0,
        stock: 0
    })
    const [selectedFile, setSelectedFile] = useState<File | null>(null)
    const [preview, setPreview] = useState<string | null>(null)
    const [variants, setVariants] = useState<VariantItem[]>([])
    const [activeCropContext, setActiveCropContext] = useState<'main' | number>('main')

    const fileInputRef = useRef<HTMLInputElement>(null)

    useEffect(() => {
        if (open) {
            if (productToEdit) {
                console.log("DEBUG: Editing Product", { list: productToEdit, detail: detailProduct })

                setFormData({
                    name: detailProduct?.name ?? productToEdit.name ?? '',
                    category_id: detailProduct?.category_id ?? productToEdit.category_id ?? 0,
                    price: detailProduct?.price ?? productToEdit.price ?? 0,
                    stock: detailProduct?.stock ?? productToEdit.stock ?? 0
                })
                setPreview(detailProduct?.image_url ?? productToEdit.image_url ?? null)

                const options = detailProduct?.options || []
                setVariants(options.map(opt => ({
                    id: opt.id,
                    name: opt.name || '',
                    additional_price: opt.additional_price || 0,
                    image_url: opt.image_url,
                    isNew: false
                })))
            } else {
                setFormData({ name: '', category_id: 0, price: 0, stock: 0 })
                setPreview(null)
                setSelectedFile(null)
                setVariants([])
            }
            setActiveCropContext('main')
        }
    }, [open, productToEdit, detailProduct])

    const handleCropSuccess = (file: File) => {
        if (activeCropContext === 'main') {
            setSelectedFile(file)
            setPreview(URL.createObjectURL(file))
        } else {
            const index = activeCropContext
            setVariants(prev => {
                const newVariants = [...prev]
                newVariants[index] = {
                    ...newVariants[index],
                    imageFile: file,
                    image_url: URL.createObjectURL(file)
                }
                return newVariants
            })
        }
    }

    const handleAddVariant = () => {
        setVariants([...variants, { name: '', additional_price: 0, isNew: true }])
    }

    const handleRemoveVariant = (index: number) => {
        if (variants[index].isNew) {
            setVariants(variants.filter((_, i) => i !== index))
        } else {
            toast.error("Cannot delete existing variants (API limitation)")
        }
    }

    const handleVariantChange = (index: number, field: keyof VariantItem, value: any) => {
        setVariants(prev => {
            const newVariants = [...prev]
            newVariants[index] = { ...newVariants[index], [field]: value }
            return newVariants
        })
    }

    const triggerVariantImageUpload = (index: number) => {
        setActiveCropContext(index)
    
        setTimeout(() => fileInputRef.current?.click(), 0)
    }

    const triggerMainImageUpload = () => {
        setActiveCropContext('main')
        setTimeout(() => fileInputRef.current?.click(), 0)
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        if (!formData.category_id) {
            toast.error("Please select a category")
            return
        }

        try {
            let productId = productToEdit?.id

            if (productToEdit && productId) {
                const payload: POSKasirInternalDtoUpdateProductRequest = {
                    name: formData.name,
                    category_id: formData.category_id,
                    price: Number(formData.price),
                    stock: Number(formData.stock)
                }
                await updateMutation.mutateAsync({ id: productId, body: payload })
            } else {
                const payload: POSKasirInternalDtoCreateProductRequest = {
                    name: formData.name,
                    category_id: formData.category_id,
                    price: Number(formData.price),
                    stock: Number(formData.stock),
                    options: []
                }
                const newProduct = await createMutation.mutateAsync(payload)
                productId = newProduct.id
            }

            if (selectedFile && productId) {
                await uploadImageMutation.mutateAsync({ id: productId, file: selectedFile })
            }

            // Handle Variants
            if (productId) {
                for (const variant of variants) {
                    let optionId = variant.id

                    if (variant.isNew) {
                        const newOpt = await createOptionMutation.mutateAsync({
                            productId,
                            body: {
                                name: variant.name,
                                additional_price: Number(variant.additional_price)
                            }
                        })
                        optionId = newOpt.id
                    } else if (optionId) {

                        await updateOptionMutation.mutateAsync({
                            productId,
                            optionId,
                            body: {
                                name: variant.name,
                                additional_price: Number(variant.additional_price)
                            }
                        })
                    }

                    if (variant.imageFile && optionId) {
                        await uploadOptionImageMutation.mutateAsync({
                            productId,
                            optionId,
                            file: variant.imageFile
                        })
                    }
                }
            }

            onOpenChange(false)
            toast.success(productToEdit ? "Product updated!" : "Product created!")

        } catch (error) {
            console.error(error)
        }
    }

    const isSubmitting = createMutation.isPending ||
        updateMutation.isPending ||
        uploadImageMutation.isPending ||
        createOptionMutation.isPending ||
        updateOptionMutation.isPending ||
        uploadOptionImageMutation.isPending

    return (
        <>
            <Dialog open={open} onOpenChange={onOpenChange}>
                <DialogContent className="sm:max-w-[500px] h-[80vh] flex flex-col p-0"> 
                    <DialogHeader className="p-6 pb-2">
                        <DialogTitle>{productToEdit ? 'Edit Product' : 'Create New Product'}</DialogTitle>
                        <DialogDescription>
                            Fill in the details below. Click save when you're done.
                        </DialogDescription>
                    </DialogHeader>

                    {/* Scrollable Content Area */}
                    <div className="flex-1 overflow-y-auto px-6 py-2">
                        {productToEdit && isLoadingDetail ? (
                            <div className="h-full flex items-center justify-center">
                                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                            </div>
                        ) : (
                            <form id="product-form" onSubmit={handleSubmit} className="grid gap-6 py-4">
                                {/* Image Upload Section */}
                                <div className="flex flex-col items-center gap-4">
                                    <Avatar
                                        className="h-24 w-24 border-2 border-muted cursor-pointer hover:opacity-80 transition-opacity"
                                        onClick={triggerMainImageUpload}>
                                        <AvatarImage src={preview || undefined} className="object-cover" />
                                        <AvatarFallback className="bg-muted">
                                            <ImageIcon className="h-8 w-8 text-muted-foreground" />
                                        </AvatarFallback>
                                    </Avatar>
                                    <Button
                                        type="button"
                                        variant="outline"
                                        size="sm"
                                        onClick={triggerMainImageUpload}
                                    >
                                        {preview ? "Change Image" : "Upload Image"}
                                    </Button>

                                    {/* Input File Hidden terhubung ke hook */}
                                    <input
                                        ref={fileInputRef}
                                        type="file"
                                        accept="image/*"
                                        className="hidden"
                                        onChange={cropper.onFileChange}
                                    />
                                </div>


                                <div className="grid gap-4">
                                    {/* Name */}
                                    <div className="grid grid-cols-4 items-center gap-4">
                                        <Label htmlFor="name" className="text-right">Name</Label>
                                        <Input
                                            id="name"
                                            value={formData.name}
                                            onChange={e => setFormData({ ...formData, name: e.target.value })}
                                            className="col-span-3"
                                            required
                                        />
                                    </div>

                                    {/* Category */}
                                    <div className="grid grid-cols-4 items-center gap-4">
                                        <Label htmlFor="category" className="text-right">Category</Label>
                                        <div className="col-span-3">
                                            <Select
                                                value={formData.category_id ? String(formData.category_id) : ""}
                                                onValueChange={val => setFormData({ ...formData, category_id: Number(val) })}
                                            >
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Select Category" />
                                                </SelectTrigger>
                                                <SelectContent>
                                                    {categories.map((cat: any) => (
                                                        <SelectItem key={cat.id} value={String(cat.id)}>
                                                            {cat.name}
                                                        </SelectItem>
                                                    ))}
                                                </SelectContent>
                                            </Select>
                                        </div>
                                    </div>

                                    {/* Price */}
                                    <div className="grid grid-cols-4 items-center gap-4">
                                        <Label htmlFor="price" className="text-right">Price</Label>
                                        <div className="col-span-3 relative">
                                            <span className="absolute left-3 top-2.5 text-sm text-muted-foreground">Rp</span>
                                            <Input
                                                id="price"
                                                type="text"
                                                inputMode="numeric"
                                                value={formData.price ? formData.price.toLocaleString('id-ID') : ''}
                                                onChange={e => {
                                                    const val = e.target.value.replace(/\D/g, '')
                                                    setFormData({ ...formData, price: Number(val) })
                                                }}
                                                className="pl-9"
                                                placeholder="0"
                                                required
                                            />
                                        </div>
                                    </div>

                                    {/* Stock */}
                                    <div className="grid grid-cols-4 items-center gap-4">
                                        <Label htmlFor="stock" className="text-right">Stock</Label>
                                        <Input
                                            id="stock"
                                            type="text"
                                            inputMode="numeric"
                                            value={formData.stock ? formData.stock.toLocaleString('id-ID') : ''}
                                            onChange={e => {
                                                const val = e.target.value.replace(/\D/g, '')
                                                setFormData({ ...formData, stock: Number(val) })
                                            }}
                                            className="col-span-3"
                                            placeholder="0"
                                            required
                                        />
                                    </div>
                                </div>


                                {/* Variants Section */}
                                <div className="space-y-4">
                                    <div className="flex items-center justify-between">
                                        <Label>Product Variants</Label>
                                        <Button type="button" variant="outline" size="sm" onClick={handleAddVariant}>
                                            <Plus className="mr-2 h-3 w-3" /> Add Variant
                                        </Button>
                                    </div>

                                    <div className="space-y-3">
                                        {variants.map((variant, index) => (
                                            <div key={index} className="flex items-start gap-3 p-3 border rounded-md">
                                                {/* Variant Image */}
                                                <div
                                                    className="h-12 w-12 shrink-0 border rounded-md cursor-pointer hover:opacity-80 flex items-center justify-center bg-muted"
                                                    onClick={() => triggerVariantImageUpload(index)}
                                                >
                                                    {variant.image_url ? (
                                                        <img src={variant.image_url} alt="" className="h-full w-full object-cover rounded-md" />
                                                    ) : (
                                                        <ImageIcon className="h-5 w-5 text-muted-foreground" />
                                                    )}
                                                </div>

                                                <div className="grid grid-cols-2 gap-2 flex-1">
                                                    <div className="space-y-1">
                                                        <Input
                                                            placeholder="Variant Name (e.g. Size XL)"
                                                            value={variant.name}
                                                            onChange={e => handleVariantChange(index, "name", e.target.value)}
                                                            required
                                                        />
                                                    </div>
                                                    <div className="space-y-1 relative">
                                                        <span className="absolute left-3 top-2.5 text-xs text-muted-foreground">Rp</span>
                                                        <Input
                                                            type="text"
                                                            inputMode="numeric"
                                                            placeholder="Additional Price"
                                                            className="pl-8"
                                                            value={variant.additional_price ? variant.additional_price.toLocaleString('id-ID') : ''}
                                                            onChange={e => {
                                                                const val = e.target.value.replace(/\D/g, '')
                                                                handleVariantChange(index, "additional_price", Number(val))
                                                            }}
                                                        />
                                                    </div>
                                                </div>

                                                <Button
                                                    type="button"
                                                    variant="ghost"
                                                    size="icon"
                                                    className="text-destructive"
                                                    onClick={() => handleRemoveVariant(index)}
                                                    disabled={!variant.isNew}
                                                    title={!variant.isNew ? "Cannot remove existing variant" : "Remove variant"}
                                                >
                                                    <Trash className="h-4 w-4" />
                                                </Button>
                                            </div>
                                        ))}
                                        {variants.length === 0 && (
                                            <p className="text-sm text-center text-muted-foreground py-2">No variants added.</p>
                                        )}
                                    </div>
                                </div> 

                              
                            </form>
                        )} 
                    </div> 

                    <DialogFooter className="p-6 pt-2">
                        <Button type="submit" form="product-form" disabled={isSubmitting || (!!productToEdit && isLoadingDetail)}>
                            {isSubmitting ? (
                                <><Loader2 className="mr-2 h-4 w-4 animate-spin" /> Saving...</>
                            ) : (
                                "Save Product"
                            )}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog >


            <ImageCropperDialog
                open={cropper.isDialogOpen}
                onOpenChange={cropper.setIsDialogOpen}
                imageSrc={cropper.imageSrc}
                crop={cropper.crop}
                zoom={cropper.zoom}
                setCrop={cropper.setCrop}
                setZoom={cropper.setZoom}
                onCropComplete={cropper.onCropComplete}
                onSave={() => cropper.onCropSave(handleCropSuccess)}
                aspectRatio={1}
            />
        </>


    )
}