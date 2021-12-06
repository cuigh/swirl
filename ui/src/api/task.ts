import ajax, { Result } from './ajax'

export interface Task {
    id: string;
    name: string;
    version: number;
    image: string;
    slot: number;
    state: string;
    serviceId: string;
    serviceName: string;
    nodeId: string;
    containerId?: string;
    pid?: number;
    exitCode?: number;
    message: string;
    error: string;
    env?: {
        name: string;
        value: string;
    }[];
    labels?: {
        name: string;
        value: string;
    }[];
    networks?: {
        id: string;
        name: string;
        ips: string[];
    }[];
    createdAt: string;
    updatedAt: string;
}

export interface SearchArgs {
    node?: string;
    service?: string;
    state?: string;
    pageIndex: number;
    pageSize: number;
}

export interface SearchResult {
    items: Task[];
    total: number;
}

export interface FindResult {
    task: Task;
    raw: string;
}

export interface FetchLogsArgs {
    id: string;
    lines: number;
    timestamps: boolean;
}

export class TaskApi {
    find(id: string) {
        return ajax.get<FindResult>('/task/find', { id })
    }

    search(args: SearchArgs) {
        return ajax.get<SearchResult>('/task/search', args)
    }

    fetchLogs(args: FetchLogsArgs) {
        return ajax.get<{
            stdout: string;
            stderr: string;
        }>('/task/fetch-logs', args)
    }
}

export default new TaskApi
