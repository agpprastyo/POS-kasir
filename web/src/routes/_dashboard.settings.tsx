import { createFileRoute, getRouteApi } from '@tanstack/react-router'
import { useState, useRef, useCallback } from 'react'
import { useSuspenseQuery } from '@tanstack/react-query' // Import hook Query
// Import meQueryOptions untuk berlangganan data user
import {
    useUpdateAvatarMutation,
    useUpdatePasswordMutation,
    meQueryOptions
} from '@/lib/api/query/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Loader2, Upload, Save, KeyRound, User, Crop as CropIcon } from 'lucide-react'
import { POSKasirInternalDtoUpdatePasswordRequest } from '@/lib/api/generated'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { Slider } from "@/components/ui/slider"
import Cropper, { Area } from 'react-easy-crop'

// Mengakses Loader Data dari parent route (_dashboard)
const routeApi = getRouteApi('/_dashboard')

export const Route = createFileRoute('/_dashboard/settings')({
    component: SettingsPage,
})

function SettingsPage() {

    const { data: user } = useSuspenseQuery(meQueryOptions())

    const profile = (user as any)?.data ?? user

    return (
        <div className="flex flex-col gap-6">
            <div>
                <h1 className="text-3xl font-bold tracking-tight">Settings</h1>
                <p className="text-muted-foreground">
                    Manage your account settings and preferences.
                </p>
            </div>

            <div className="grid gap-6 md:grid-cols-2">
                {/* Kartu 1: Update Profile Picture */}
                <UpdateAvatarCard
                    currentAvatar={profile?.avatar}
                    username={profile?.username}
                />

                {/* Kartu 2: Update Password */}
                <UpdatePasswordCard />
            </div>
        </div>
    )
}

function UpdateAvatarCard({ currentAvatar, username }: { currentAvatar?: string, username?: string }) {
    const [preview, setPreview] = useState<string | null>(null)
    const [selectedFile, setSelectedFile] = useState<File | null>(null)
    const [message, setMessage] = useState<{ type: 'success' | 'error', text: string } | null>(null)


    const [imageSrc, setImageSrc] = useState<string | null>(null)
    const [crop, setCrop] = useState({ x: 0, y: 0 })
    const [zoom, setZoom] = useState(1)
    const [croppedAreaPixels, setCroppedAreaPixels] = useState<Area | null>(null)
    const [isCropDialogOpen, setIsCropDialogOpen] = useState(false)

    const fileInputRef = useRef<HTMLInputElement>(null)
    const mutation = useUpdateAvatarMutation()

    const onCropComplete = useCallback((_croppedArea: Area, croppedAreaPixels: Area) => {
        setCroppedAreaPixels(croppedAreaPixels)
    }, [])

    const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files.length > 0) {
            const file = e.target.files[0]
            const imageDataUrl = await readFile(file)
            setImageSrc(imageDataUrl as string)
            setIsCropDialogOpen(true)
            e.target.value = ''
        }
    }

    const handleCropSave = async () => {
        if (!imageSrc || !croppedAreaPixels) return

        try {
            const croppedImageBlob = await getCroppedImg(imageSrc, croppedAreaPixels)
            const file = new File([croppedImageBlob], "avatar.jpg", { type: "image/jpeg" })

            setSelectedFile(file)
            setPreview(URL.createObjectURL(file))
            setMessage(null)
            setIsCropDialogOpen(false)
        } catch (e) {
            console.error(e)
            setMessage({ type: 'error', text: 'Failed to crop image.' })
        }
    }

    const handleSave = async () => {
        if (!selectedFile) return

        try {
            await mutation.mutateAsync(selectedFile)
            setMessage({ type: 'success', text: 'Profile picture updated successfully!' })

            setPreview(null)
            setSelectedFile(null)
            setImageSrc(null)
            if (fileInputRef.current) fileInputRef.current.value = ''

        } catch (error: any) {
            setMessage({ type: 'error', text: 'Failed to update avatar.' })
        }
    }

    return (
        <>
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <User className="h-5 w-5" /> Profile Picture
                    </CardTitle>
                    <CardDescription>
                        Update your profile picture.
                    </CardDescription>
                </CardHeader>
                <CardContent className="flex flex-col items-center gap-6">
                    <Avatar className="h-32 w-32 border-4 border-muted">

                        <AvatarImage src={preview || currentAvatar || "https://github.com/shadcn.png"} />
                        <AvatarFallback className="text-4xl">
                            {username?.slice(0, 2).toUpperCase() ?? 'US'}
                        </AvatarFallback>
                    </Avatar>

                    <div className="flex w-full max-w-sm items-center gap-2">
                        <Input
                            ref={fileInputRef}
                            type="file"
                            accept="image/*"
                            onChange={handleFileChange}
                            className="cursor-pointer"
                        />
                    </div>

                    {message && (
                        <Alert
                            variant={(message.type === 'error' ? 'destructive' : 'default') as "default" | "destructive"}
                            className={message.type === 'success' ? 'border-green-500 text-green-500' : ''}
                        >
                            <AlertDescription>{message.text}</AlertDescription>
                        </Alert>
                    )}
                </CardContent>
                <CardFooter className="justify-end border-t bg-muted/20 px-6 py-4">
                    <Button
                        onClick={handleSave}
                        disabled={!selectedFile || mutation.isPending}
                    >
                        {mutation.isPending ? (
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        ) : (
                            <Upload className="mr-2 h-4 w-4" />
                        )}
                        Upload New Picture
                    </Button>
                </CardFooter>
            </Card>

            {/* Dialog Crop */}
            <Dialog open={isCropDialogOpen} onOpenChange={setIsCropDialogOpen}>
                <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle>Crop Image</DialogTitle>
                        <DialogDescription>
                            Adjust the image to fit the square aspect ratio.
                        </DialogDescription>
                    </DialogHeader>

                    <div className="relative h-[300px] w-full overflow-hidden rounded-md border bg-slate-900">
                        {imageSrc && (
                            <Cropper
                                image={imageSrc}
                                crop={crop}
                                zoom={zoom}
                                aspect={1} // Force Square Ratio
                                onCropChange={setCrop}
                                onCropComplete={onCropComplete}
                                onZoomChange={setZoom}
                            />
                        )}
                    </div>

                    <div className="flex items-center space-x-2 pt-2">
                        <span className="text-xs text-muted-foreground">Zoom</span>
                        <Slider
                            value={[zoom]}
                            min={1}
                            max={3}
                            step={0.1}
                            onValueChange={(vals) => setZoom(vals[0])}
                            className="flex-1"
                        />
                    </div>

                    <DialogFooter className="sm:justify-between">
                        <Button
                            variant="secondary"
                            onClick={() => {
                                setIsCropDialogOpen(false)
                                setImageSrc(null)
                            }}
                        >
                            Cancel
                        </Button>
                        <Button onClick={handleCropSave}>
                            <CropIcon className="mr-2 h-4 w-4" />
                            Crop & Save
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </>
    )
}

// --- KOMPONEN: Update Password ---
function UpdatePasswordCard() {
    const [formData, setFormData] = useState({
        old_password: '',
        new_password: '',
        confirm_password: ''
    })
    const [message, setMessage] = useState<{ type: 'success' | 'error', text: string } | null>(null)

    const mutation = useUpdatePasswordMutation()

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target
        setFormData(prev => ({ ...prev, [name]: value }))
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setMessage(null)

        if (formData.new_password !== formData.confirm_password) {
            setMessage({ type: 'error', text: 'New passwords do not match.' })
            return
        }

        try {
            const payload: POSKasirInternalDtoUpdatePasswordRequest = {
                old_password: formData.old_password,
                new_password: formData.new_password
            }

            await mutation.mutateAsync(payload)
            setMessage({ type: 'success', text: 'Password updated successfully!' })
            setFormData({ old_password: '', new_password: '', confirm_password: '' })
        } catch (error: any) {
            const msg = error?.response?.data?.message ?? 'Failed to update password.'
            setMessage({ type: 'error', text: msg })
        }
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <KeyRound className="h-5 w-5" /> Security
                </CardTitle>
                <CardDescription>
                    Change your password to keep your account secure.
                </CardDescription>
            </CardHeader>
            <form onSubmit={handleSubmit}>
                <CardContent className="grid gap-4">
                    {message && (
                        <Alert
                            variant={(message.type === 'error' ? 'destructive' : 'default') as "default" | "destructive"}
                            className={message.type === 'success' ? 'border-green-500 text-green-500' : ''}
                        >
                            <AlertDescription>{message.text}</AlertDescription>
                        </Alert>
                    )}

                    <div className="grid gap-2">
                        <Label htmlFor="old_password">Current Password</Label>
                        <Input
                            id="old_password"
                            name="old_password"
                            type="password"
                            value={formData.old_password}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div className="grid gap-2">
                        <Label htmlFor="new_password">New Password</Label>
                        <Input
                            id="new_password"
                            name="new_password"
                            type="password"
                            value={formData.new_password}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div className="grid gap-2">
                        <Label htmlFor="confirm_password">Confirm New Password</Label>
                        <Input
                            id="confirm_password"
                            name="confirm_password"
                            type="password"
                            value={formData.confirm_password}
                            onChange={handleChange}
                            required
                        />
                    </div>
                </CardContent>
                <CardFooter className="justify-end border-t bg-muted/20 px-6 py-4">
                    <Button type="submit" disabled={mutation.isPending}>
                        {mutation.isPending ? (
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        ) : (
                            <Save className="mr-2 h-4 w-4" />
                        )}
                        Change Password
                    </Button>
                </CardFooter>
            </form>
        </Card>
    )
}

// --- UTILITIES ---

function readFile(file: File) {
    return new Promise((resolve) => {
        const reader = new FileReader()
        reader.addEventListener('load', () => resolve(reader.result), false)
        reader.readAsDataURL(file)
    })
}

const createImage = (url: string): Promise<HTMLImageElement> =>
    new Promise((resolve, reject) => {
        const image = new Image()
        image.addEventListener('load', () => resolve(image))
        image.addEventListener('error', (error) => reject(error))
        image.setAttribute('crossOrigin', 'anonymous')
        image.src = url
    })

async function getCroppedImg(imageSrc: string, pixelCrop: Area): Promise<Blob> {
    const image = await createImage(imageSrc)
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')

    if (!ctx) {
        throw new Error('No 2d context')
    }

    canvas.width = pixelCrop.width
    canvas.height = pixelCrop.height

    ctx.drawImage(
        image,
        pixelCrop.x,
        pixelCrop.y,
        pixelCrop.width,
        pixelCrop.height,
        0,
        0,
        pixelCrop.width,
        pixelCrop.height
    )

    return new Promise((resolve, reject) => {
        canvas.toBlob((blob) => {
            if (!blob) {
                reject(new Error('Canvas is empty'))
                return
            }
            resolve(blob)
        }, 'image/jpeg')
    })
}