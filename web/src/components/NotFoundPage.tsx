import { Link } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { FileQuestion } from 'lucide-react'

export function NotFoundPage() {
    return (
        <div className="flex h-screen w-full flex-col items-center justify-center gap-4 bg-background px-4 text-center">
            <div className="rounded-full bg-muted p-4">
                <FileQuestion className="h-10 w-10 text-muted-foreground" />
            </div>
            <div className="space-y-2">
                <h1 className="text-4xl font-bold tracking-tighter sm:text-5xl">404</h1>
                <h2 className="text-2xl font-semibold tracking-tight">Page Not Found</h2>
                <p className="max-w-[500px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
                    Sorry, the page you are looking for does not exist or has been moved.
                </p>
            </div>
            <Button asChild variant="default" size="lg">
                <Link to="/">Go Back Home</Link>
            </Button>
        </div>
    )
}
