import { activityLogsApi } from "../client";
import { POSKasirInternalDtoActivityLogListResponse } from "@/lib/api/generated";
import { queryOptions, useQuery } from "@tanstack/react-query";
import { z } from "zod";

export const activityLogsSearchSchema = z.object({
    page: z.number().min(1).catch(1),
    limit: z.number().min(1).max(100).catch(10),
    search: z.string().optional(),
    start_date: z.string().optional(),
    end_date: z.string().optional(),
    user_id: z.string().optional(),
    entity_type: z.string().optional(),
    action_type: z.string().optional(),
});

export type ActivityLogsSearch = z.infer<typeof activityLogsSearchSchema>;

export const activityLogsListQueryOptions = (search: ActivityLogsSearch) => queryOptions({
    queryKey: ['activity-logs', search],
    queryFn: async () => {
        const response = await activityLogsApi.activityLogsGet(
            search.page,
            search.limit,
            search.search,
            search.start_date,
            search.end_date,
            search.user_id,
            search.entity_type,
            search.action_type
        );
        return (response.data as any).data as POSKasirInternalDtoActivityLogListResponse;
    }
});

export const useActivityLogsList = (search: ActivityLogsSearch) => {
    return useQuery(activityLogsListQueryOptions(search));
};
