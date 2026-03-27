import { activityLogsApi } from "../client";
import {
    InternalActivitylogActivityLogListResponse,
    ActivityLogsGetEntityTypeEnum,
    ActivityLogsGetActionTypeEnum
} from "@/lib/api/generated";
import { queryOptions, useQuery } from "@tanstack/react-query";
import { z } from "zod";
import { useRBAC } from '@/lib/auth/rbac';

export const activityLogsSearchSchema = z.object({
    page: z.number().min(1).catch(1),
    limit: z.number().min(1).max(100).catch(10),
    search: z.string().optional(),
    start_date: z.string().optional(),
    end_date: z.string().optional(),
    user_id: z.string().optional(),
    entity_type: z.preprocess((val) => typeof val === 'string' ? val.toUpperCase() : val, z.enum(Object.values(ActivityLogsGetEntityTypeEnum)).optional().catch(undefined)).optional(),
    action_type: z.preprocess((val) => {
        if (typeof val !== 'string') return val;
        const upper = val.toUpperCase();
        if (upper === 'LOGIN') return ActivityLogsGetActionTypeEnum.LoginSuccess;
        if (upper === 'LOGOUT') return undefined;
        return upper;
    }, z.enum(Object.values(ActivityLogsGetActionTypeEnum)).optional().catch(undefined)).optional(),
});

export type ActivityLogsSearch = {
    page: number;
    limit: number;
    search?: string;
    start_date?: string;
    end_date?: string;
    user_id?: string;
    entity_type?: ActivityLogsGetEntityTypeEnum;
    action_type?: ActivityLogsGetActionTypeEnum;
};

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
        return (response.data as any).data as InternalActivitylogActivityLogListResponse;
    }
});

export const useActivityLogsList = (search: ActivityLogsSearch) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/activity-logs');
    const query = useQuery({
        ...activityLogsListQueryOptions(search),
        enabled: isAllowed
    });
    return { ...query, isAllowed };
};
