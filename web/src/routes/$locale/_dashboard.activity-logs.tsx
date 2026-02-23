import { Button } from '@/components/ui/button';
import { Card, CardContent, CardTitle, CardHeader } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import {
    POSKasirInternalUserRepositoryUserRole,
    InternalActivitylogActivityLogResponse,
    ActivityLogsGetActionTypeEnum,
    ActivityLogsGetEntityTypeEnum
} from '@/lib/api/generated';
import { activityLogsListQueryOptions, activityLogsSearchSchema, ActivityLogsSearch } from '@/lib/api/query/activity-logs';
import { meQueryOptions } from '@/lib/api/query/auth';
import { cn } from '@/lib/utils';
import { createFileRoute, redirect, useNavigate } from '@tanstack/react-router';
import { useSuspenseQuery } from '@tanstack/react-query';
import { format } from 'date-fns';
import { CalendarIcon, ChevronLeft, ChevronRight, Search, X, ChevronsUpDown } from 'lucide-react';
import { Calendar } from '@/components/ui/calendar';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';
import { useTranslation } from 'react-i18next';

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
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">{t('activity_logs.title')}</h2>
                    <p className="text-muted-foreground">
                        {t('activity_logs.description')}
                    </p>
                </div>
                <Button variant="outline" onClick={() => navigate({ search: { page: 1, limit: 10 } })}>
                    <X className="mr-2 h-4 w-4" />
                    {t('common.clear_filters')}
                </Button>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>{t('activity_logs.list_title')}</CardTitle>
                    <div className="flex flex-wrap items-center gap-4 py-4">
                        <div className="relative flex-1 min-w-[200px]">
                            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                            <Input
                                placeholder={t('activity_logs.search_placeholder')}
                                className="pl-8"
                                value={search.search || ''}
                                onChange={(e) => handleSearch(e.target.value)}
                            />
                        </div>

                        <div className="w-[180px]">
                            <Select value={search.action_type || 'all'} onValueChange={(val) => updateFilter('action_type', val === 'all' ? undefined : val)}>
                                <SelectTrigger>
                                    <SelectValue placeholder={t('activity_logs.filters.action_type')} />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="all">{t('activity_logs.filters.all_actions')}</SelectItem>
                                    {Object.values(ActivityLogsGetActionTypeEnum).map((action) => (
                                        <SelectItem key={action} value={action}>
                                            {action.replace(/_/g, ' ')}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="w-[180px]">
                            <Select value={search.entity_type || 'all'} onValueChange={(val) => updateFilter('entity_type', val === 'all' ? undefined : val)}>
                                <SelectTrigger>
                                    <SelectValue placeholder={t('activity_logs.filters.entity_type')} />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="all">{t('activity_logs.filters.all_entities')}</SelectItem>
                                    {Object.values(ActivityLogsGetEntityTypeEnum).map((entity) => (
                                        <SelectItem key={entity} value={entity}>
                                            {entity}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="w-[180px]">
                            <Input
                                placeholder={t('activity_logs.filters.user_id')}
                                value={search.user_id || ''}
                                onChange={(e) => updateFilter('user_id', e.target.value || undefined)}
                            />
                        </div>

                        <Popover>
                            <PopoverTrigger asChild>
                                <Button
                                    variant={"outline"}
                                    className={cn(
                                        "w-[240px] justify-start text-left font-normal",
                                        !search.start_date && "text-muted-foreground"
                                    )}
                                >
                                    <CalendarIcon className="mr-2 h-4 w-4" />
                                    {search.start_date ? (
                                        search.end_date ? (
                                            <>
                                                {format(new Date(search.start_date), "LLL dd, y")} -{" "}
                                                {format(new Date(search.end_date), "LLL dd, y")}
                                            </>
                                        ) : (
                                            format(new Date(search.start_date), "LLL dd, y")
                                        )
                                    ) : (
                                        <span>{t('activity_logs.filters.date_range')}</span>
                                    )}
                                </Button>
                            </PopoverTrigger>
                            <PopoverContent className="w-auto p-0" align="start">
                                <Calendar
                                    initialFocus
                                    mode="range"
                                    defaultMonth={search.start_date ? new Date(search.start_date) : undefined}
                                    selected={{
                                        from: search.start_date ? new Date(search.start_date) : undefined,
                                        to: search.end_date ? new Date(search.end_date) : undefined,
                                    }}
                                    onSelect={(range) => {
                                        navigate({
                                            search: (prev) => ({
                                                ...prev,
                                                start_date: range?.from ? format(range.from, 'yyyy-MM-dd') : undefined,
                                                end_date: range?.to ? format(range.to, 'yyyy-MM-dd') : undefined,
                                                page: 1
                                            }),
                                        });
                                    }}
                                    numberOfMonths={2}
                                />
                            </PopoverContent>
                        </Popover>
                    </div>
                </CardHeader>
                <CardContent>
                    <div className="rounded-md border">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>{t('activity_logs.columns.date')}</TableHead>
                                    <TableHead>{t('activity_logs.columns.user')}</TableHead>
                                    <TableHead>{t('activity_logs.columns.action')}</TableHead>
                                    <TableHead>{t('activity_logs.columns.entity')}</TableHead>
                                    <TableHead>{t('activity_logs.columns.details')}</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {data?.logs?.map((log: InternalActivitylogActivityLogResponse) => (
                                    <TableRow key={log.id}>
                                        <TableCell>{format(new Date(log.created_at!), 'PPP p')}</TableCell>
                                        <TableCell>
                                            <div className="font-medium">{log.user_name}</div>
                                            <div className="text-xs text-muted-foreground">{t('activity_logs.table.id_prefix')} {log.user_id?.substring(0, 8)}...</div>
                                        </TableCell>
                                        <TableCell>
                                            <span className="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 border-transparent bg-primary text-primary-foreground hover:bg-primary/80">
                                                {log.action_type}
                                            </span>
                                        </TableCell>
                                        <TableCell>
                                            <div className="font-medium">{log.entity_type}</div>
                                            <div className="text-xs text-muted-foreground">{log.entity_id}</div>
                                        </TableCell>
                                        <TableCell>
                                            <Collapsible>
                                                <div className="flex items-center gap-2">
                                                    <CollapsibleTrigger asChild>
                                                        <Button variant="ghost" size="sm" className="h-6 w-6 p-0">
                                                            <ChevronsUpDown className="h-4 w-4" />
                                                            <span className="sr-only">{t('activity_logs.table.toggle')}</span>
                                                        </Button>
                                                    </CollapsibleTrigger>
                                                    <span className="text-xs text-muted-foreground truncate max-w-[200px]">
                                                        {JSON.stringify(log.details)}
                                                    </span>
                                                </div>
                                                <CollapsibleContent>
                                                    <pre className="mt-2 w-[300px] overflow-auto rounded-md bg-muted p-2 text-xs">
                                                        {JSON.stringify(log.details, null, 2)}
                                                    </pre>
                                                </CollapsibleContent>
                                            </Collapsible>
                                        </TableCell>
                                    </TableRow>
                                ))}
                                {data?.logs?.length === 0 && (
                                    <TableRow>
                                        <TableCell colSpan={5} className="h-24 text-center">
                                            {t('common.no_results')}
                                        </TableCell>
                                    </TableRow>
                                )}
                            </TableBody>
                        </Table>
                    </div>

                    <div className="flex items-center justify-end space-x-2 py-4">
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handlePageChange(search.page - 1)}
                            disabled={search.page <= 1}
                        >
                            <ChevronLeft className="h-4 w-4" />
                            {t('common.previous')}
                        </Button>
                        <div className="text-sm font-medium">
                            {t('common.page_info', {
                                current: data?.page,
                                total: data?.total_pages || 1
                            })}
                        </div>
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handlePageChange(search.page + 1)}
                            disabled={search.page >= (data?.total_pages || 1)}
                        >
                            {t('common.next')}
                            <ChevronRight className="h-4 w-4" />
                        </Button>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
