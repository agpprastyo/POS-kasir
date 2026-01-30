import { Button } from "@/components/ui/button.tsx";
import { useTranslation } from 'react-i18next';
import {POSKasirInternalCommonPaginationPagination} from "@/lib/api/generated";

interface NewPaginationProps {
    pagination: POSKasirInternalCommonPaginationPagination | undefined;
    onClickPrev: () => void;
    onClickNext: () => void;
}

export function NewPagination({ pagination, onClickPrev, onClickNext }: NewPaginationProps) {
    const { t } = useTranslation();

    if (!pagination) {
        return null;
    }

    return (
        <div className="flex items-center justify-end gap-2">
            <div className="text-sm text-muted-foreground">
                {t('pagination.page_info', { current: pagination.current_page, total: pagination.total_page })}
            </div>
            <Button
                variant="outline"
                size="sm"
                onClick={onClickPrev}
                disabled={(pagination.current_page || 1) <= 1}
            >
                {t('pagination.previous')}
            </Button>
            <Button
                variant="outline"
                size="sm"
                onClick={onClickNext}
                disabled={(pagination.current_page || 1) >= (pagination.total_page || 1)}
            >
                {t('pagination.next')}
            </Button>
        </div>
    );
}