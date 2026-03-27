import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
    POSKasirInternalUserRepositoryUserRole,
} from '@/lib/api/generated';
import { activityLogsListQueryOptions, activityLogsSearchSchema, ActivityLogsSearch } from '@/lib/api/query/activity-logs';
import { meQueryOptions } from '@/lib/api/query/auth';
import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router';
import { useSuspenseQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { ActivityLogsHeader } from "@/components/activity-logs/ActivityLogsHeader"
import { ActivityLogsFilters } from "@/components/activity-logs/ActivityLogsFilters"
import { ActivityLogsTable } from "@/components/activity-logs/ActivityLogsTable"

export const Route = createFileRoute('/$locale/_dashboard/activity-logs')(({
    validateSearch: activityLogsSearchSchema,
    loaderDeps: ({ search }: any) => search,
    loader: ({ context: { queryClient }, deps }: any) => {
        queryClient.ensureQueryData(activityLogsListQueryOptions(deps));
    },
    beforeLoad: async ({ context: { queryClient } }: any) => {
        const user = await queryClient.ensureQueryData(meQueryOptions());
        const allowedRoles: POSKasirInternalUserRepositoryUserRole[] = [
            POSKasirInternalUserRepositoryUserRole.UserRoleAdmin,
        ];
        if (!user || !user.role || !allowedRoles.includes(user.role)) {
            throw redirect({
                to: '/$locale/login',
                params: { locale: 'en' },
                search: {
                    error: 'Unauthorized',
                },
            });
        }
    },
    component: ActivityLogsPage,
} as any));

function ActivityLogsPage() {
    const { t } = useTranslation();
    const search: ActivityLogsSearch = Route.useSearch();
    const navigate = useNavigate({ from: Route.fullPath });
    const { data } = useSuspenseQuery(activityLogsListQueryOptions(search));

    const handleSearch = (term: string) => {
        navigate({
            search: (prev) => ({ ...prev, search: term, page: 1 }),
        });
    };

    const handlePageChange = (newPage: number) => {
        navigate({
            search: (prev) => ({ ...prev, page: newPage }),
        });
    };

    const updateFilter = (key: string, value: string | undefined) => {
        navigate({
            search: (prev) => ({ ...prev, [key]: value, page: 1 }),
        });
    };

    return (
        <div className="space-y-6">
            <ActivityLogsHeader 
                t={t}
                onClearFilters={() => navigate({ search: { page: 1, limit: 10 } })}
            />

            <Card className="border-0 shadow-sm">
                <CardHeader>
                    <CardTitle>{t('activity_logs.list_title')}</CardTitle>
                    <ActivityLogsFilters 
                        t={t}
                        search={search}
                        onSearch={handleSearch}
                        updateFilter={updateFilter}
                        onDateChange={({ from, to }) => {
                            navigate({
                                search: (prev) => ({
                                    ...prev,
                                    start_date: from || undefined,
                                    end_date: to || undefined,
                                    page: 1
                                }),
                            });
                        }}
                    />
                </CardHeader>
                <CardContent>
                    <ActivityLogsTable 
                        t={t}
                        logs={data?.logs}
                        page={search.page}
                        totalPages={data?.total_pages || 1}
                        onPageChange={handlePageChange}
                    />
                </CardContent>
            </Card>
        </div>
    );
}
