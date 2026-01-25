import { Button } from "@/components/ui/button.tsx";
import { useTranslation } from 'react-i18next';

export function NewPagination(props: { pagination: any, onClick: () => void, onClick1: () => void }) {
    const { t } = useTranslation();
    return <>
        {props.pagination && (
            <div className="flex items-center justify-end gap-2">
                <div className="text-sm text-muted-foreground">
                    {t('pagination.page_info', { current: props.pagination.current_page, total: props.pagination.total_page })}
                </div>
                <Button
                    variant="outline"
                    size="sm"
                    onClick={props.onClick}
                    disabled={(props.pagination.current_page || 1) <= 1}
                >
                    {t('pagination.previous')}
                </Button>
                <Button
                    variant="outline"
                    size="sm"
                    onClick={props.onClick1}
                    disabled={(props.pagination.current_page || 1) >= (props.pagination.total_page || 1)}
                >
                    {t('pagination.next')}
                </Button>
            </div>
        )}
    </>;
}