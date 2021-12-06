import ajax, { Result } from './ajax'

export interface Service {
    id: string;
    name: string;
    image: string;
    version: number;
    mode: string;
    command: string;
    args: string;
    dir: string;
    user: string;
    replicas: number;
    runningTasks: number;
    desiredTasks: number;
    completedTasks: number;
    networks: string[],
    labels?: {
        name: string;
        value: string;
    }[];
    containerLabels?: {
        name: string;
        value: string;
    }[];
    env?: {
        name: string;
        value: string;
    }[];
    update?: {
        state?: string;
        message?: string;
    };
    endpoint: {
        mode: string;
        ports: {
            name: string;
            protocol: string;
            mode: string;
            targetPort: string;
            publishedPort: string;
        }[],
        vips: {
            id: string;
            name: string;
            ip: string;
        }[],
    },
    mounts?: {
        type: string;
        source: string;
        target: string;
        readonly: boolean;
        consistency: string;
    }[];
    resource: {
        limit: {
            cpu: number;
            memory: string;
        };
        reserve: {
            cpu: number;
            memory: string;
        };
    };
    placement: {
        constraints: {
            name: string;
            value: string;
            op: string;
        }[];
        preferences: string[];
    };
    configs?: {
        key: string;
        path: string;
        uid: string;
        gid: string;
        mode: number;
    }[];
    secrets?: {
        key: string;
        path: string;
        uid: string;
        gid: string;
        mode: number;
    }[];
    updatePolicy: {
        parallelism: number;
        delay: string;
        failureAction: string;
        order: string;
    };
    rollbackPolicy: {
        parallelism: number;
        delay: string;
        failureAction: string;
        order: string;
    };
    restartPolicy: {
        condition: string;
        maxAttempts: number;
        delay: string;
        window: string;
    };
    logDriver: {
        name: string;
        options?: {
            name: string;
            value: string;
        }[];
    }
    dns: {
        servers: string[];
        search: string[];
        options: string[];
    }
    hosts: string[];
    hostname: string;
    createdAt: string;
    updatedAt: string;
}

export interface SearchArgs {
    name?: string;
    mode?: string;
    pageIndex: number;
    pageSize: number;
}

export interface SearchResult {
    items: Service[];
    total: number;
}

export interface FindResult {
    service: Service;
    raw: string;
}

export interface FetchLogsArgs {
    id: string;
    lines: number;
    timestamps: boolean;
}

export class ServiceApi {
    find(name: string, status: boolean = false) {
        return ajax.get<FindResult>('/service/find', { name, status })
    }

    search(args: SearchArgs) {
        return ajax.get<SearchResult>('/service/search', args)
    }

    save(service: Service) {
        return ajax.post<Result<Object>>('/service/save', service)
    }

    delete(name: string) {
        return ajax.post<Result<Object>>('/service/delete', { name })
    }

    restart(name: string) {
        return ajax.post<Result<Object>>('/service/restart', { name })
    }

    rollback(name: string) {
        return ajax.post<Result<Object>>('/service/rollback', { name })
    }

    scale(name: string, count: number, version: number) {
        return ajax.post<Result<Object>>('/service/scale', { name, count, version })
    }

    fetchLogs(args: FetchLogsArgs) {
        return ajax.get<{
            stdout: string;
            stderr: string;
        }>('/service/fetch-logs', args)
    }
}

export default new ServiceApi
