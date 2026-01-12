import {useCallback, useState} from 'react'
import {Area} from 'react-easy-crop'
import {getCroppedImg, readFile} from '@/lib/utils'

export function useImageCropper() {
    const [imageSrc, setImageSrc] = useState<string | null>(null)
    const [crop, setCrop] = useState({x: 0, y: 0})
    const [zoom, setZoom] = useState(1)
    const [croppedAreaPixels, setCroppedAreaPixels] = useState<Area | null>(null)
    const [isDialogOpen, setIsDialogOpen] = useState(false)


    const onFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files.length > 0) {
            const file = e.target.files[0]
            const imageDataUrl = await readFile(file)
            setImageSrc(imageDataUrl as string)
            setIsDialogOpen(true)

            e.target.value = ''
        }
    }

    const onCropComplete = useCallback((_croppedArea: Area, croppedAreaPixels: Area) => {
        setCroppedAreaPixels(croppedAreaPixels)
    }, [])


    const onCropSave = async (onSuccess: (file: File) => void) => {
        if (!imageSrc || !croppedAreaPixels) return

        try {
            const croppedImageBlob = await getCroppedImg(imageSrc, croppedAreaPixels)
            const file = new File([croppedImageBlob], "cropped-image.jpg", {type: "image/jpeg"})

            onSuccess(file)
            handleClose()
        } catch (e) {
            console.error(e)
        }
    }

    const handleClose = () => {
        setIsDialogOpen(false)
        setImageSrc(null)
        setZoom(1)
        setCrop({x: 0, y: 0})
    }

    return {
        imageSrc,
        crop,
        zoom,
        isDialogOpen,
        setCrop,
        setZoom,
        setIsDialogOpen,
        onFileChange,
        onCropComplete,
        onCropSave,
        handleClose
    }
}