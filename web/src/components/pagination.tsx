import {Button} from "@/components/ui/button.tsx";

export function NewPagination(props: { pagination: any, onClick: () => void, onClick1: () => void }) {
    return <>
        {props.pagination && (
            <div className="flex items-center justify-end gap-2">
                <div className="text-sm text-muted-foreground">
                    Page {props.pagination.current_page} of {props.pagination.total_page}
                </div>
                <Button
                    variant="outline"
                    size="sm"
                    onClick={props.onClick}
                    disabled={(props.pagination.current_page || 1) <= 1}
                >
                    Previous
                </Button>
                <Button
                    variant="outline"
                    size="sm"
                    onClick={props.onClick1}
                    disabled={(props.pagination.current_page || 1) >= (props.pagination.total_page || 1)}
                >
                    Next
                </Button>
            </div>
        )}
    </>;
}