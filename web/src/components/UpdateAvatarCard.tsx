import {ReactNode, useRef, useState} from "react";
import {useUpdateAvatarMutation} from "@/lib/api/query/auth.ts";
import {useImageCropper} from "@/hooks/use-image-cropper";
import {ImageCropperDialog} from "@/components/common/ImageCropperDialog";
import {Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import {Loader2, Upload, User} from "lucide-react";
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Alert, AlertDescription} from "@/components/ui/alert.tsx";
import {Button} from "@/components/ui/button.tsx";

export function UpdateAvatarCard({currentAvatar, username}: { currentAvatar?: string, username?: string }) {
    const [preview, setPreview] = useState<string | null>(null)
    const [selectedFile, setSelectedFile] = useState<File | null>(null)
    const [message, setMessage] = useState<{ type: 'success' | 'error', text: string } | null>(null)

    const cropper = useImageCropper()
    const fileInputRef = useRef<HTMLInputElement>(null)
    const mutation = useUpdateAvatarMutation()

    const handleCropSuccess = (file: File) => {
        setSelectedFile(file)
        setPreview(URL.createObjectURL(file))
        setMessage(null)
    }

    const handleSave = async () => {
        if (!selectedFile) return
        try {
            await mutation.mutateAsync(selectedFile)
            setMessage({type: 'success', text: 'Profile picture updated successfully!'})
            setPreview(null)
            setSelectedFile(null)
            if (fileInputRef.current) fileInputRef.current.value = ''
        } catch (error: any) {
            setMessage({type: 'error', text: 'Failed to update avatar.'})
        }
    }

    return (
        <>
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <User className="h-5 w-5"/> Profile Picture
                    </CardTitle>
                    <CardDescription>
                        Update your profile picture.
                    </CardDescription>
                </CardHeader>
                <CardContent className="flex flex-col items-center gap-6">
                    <Avatar className="h-32 w-32 border-4 border-muted">
                        <AvatarImage src={preview || currentAvatar || "https://github.com/shadcn.png"}/>
                        <AvatarFallback className="text-4xl">
                            {username?.slice(0, 2).toUpperCase() ?? 'US'}
                        </AvatarFallback>
                    </Avatar>

                    <div className="flex w-full max-w-sm items-center gap-2">
                        <Input
                            ref={fileInputRef}
                            type="file"
                            accept="image/*"
                            onChange={cropper.onFileChange}
                            className="cursor-pointer"
                        />
                    </div>

                    {message && (
                        <Alert
                            variant={(message.type === 'error' ? 'destructive' : 'default') as "default" | "destructive"}
                            className={message.type === 'success' ? 'border-green-500 text-green-500' : ''}
                        >
                            <AlertDescription>{message.text}</AlertDescription>
                        </Alert> as ReactNode
                    )}
                </CardContent>
                <CardFooter className="justify-end border-t bg-muted/20 px-6 py-4">
                    <Button onClick={handleSave} disabled={!selectedFile || mutation.isPending}>
                        {mutation.isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin"/> :
                            <Upload className="mr-2 h-4 w-4"/>}
                        Upload New Picture
                    </Button>
                </CardFooter>
            </Card>

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