import Cropper, {Area} from "react-easy-crop"
import {Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog"
import {Button} from "@/components/ui/button"
import {Slider} from "@/components/ui/slider"
import {Crop as CropIcon} from "lucide-react"

interface ImageCropperDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    imageSrc: string | null
    crop: { x: number; y: number }
    zoom: number
    setCrop: (crop: { x: number; y: number }) => void
    setZoom: (zoom: number) => void
    onCropComplete: (croppedArea: Area, croppedAreaPixels: Area) => void
    onSave: () => void
    aspectRatio?: number
}

export function ImageCropperDialog({
                                       open,
                                       onOpenChange,
                                       imageSrc,
                                       crop,
                                       zoom,
                                       setCrop,
                                       setZoom,
                                       onCropComplete,
                                       onSave,
                                       aspectRatio = 1
                                   }: ImageCropperDialogProps) {
    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-md">
                <DialogHeader>
                    <DialogTitle>Crop Image</DialogTitle>
                    <DialogDescription>
                        Adjust the image to fit the required aspect ratio.
                    </DialogDescription>
                </DialogHeader>

                <div className="relative h-[300px] w-full overflow-hidden rounded-md border bg-slate-900">
                    {imageSrc && (
                        <Cropper
                            image={imageSrc}
                            crop={crop}
                            zoom={zoom}
                            aspect={aspectRatio as number}
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
                        onClick={() => onOpenChange(false)}
                        type="button"
                    >
                        Cancel
                    </Button>
                    <Button onClick={onSave} type="button">
                        <CropIcon className="mr-2 h-4 w-4"/>
                        Crop & Save
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}