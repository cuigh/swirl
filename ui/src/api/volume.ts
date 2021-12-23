import ajax, { Result } from './ajax'

export interface Volume {
    name: string;
    driver: string;
    customDriver: string;
    mountPoint: string;
    createdAt: string;
    scope: string;
    refCount: number;
    size: number;
    labels?: {
        name: string;
        value: string;
    }[];
    options?: {
        name: string;
        value: string;
    }[];
}

export interface SearchArgs {
    node?: string;
    name?: string;
    pageIndex: number;
    pageSize: number;
}

export interface SearchResult {
    items: Volume[];
    total: number;
}

export interface FindResult {
    volume: Volume;
    raw: string;
}

export class VolumeApi {
    find(node: string, name: string) {
        return ajax.get<FindResult>('/volume/find', { node, name })
    }

    search(args: SearchArgs) {
        return ajax.get<SearchResult>('/volume/search', args)
    }

    delete(node: string, name: string) {
        return ajax.post<Result<Object>>('/volume/delete', { node, name })
    }

    save(v: Volume) {
        return ajax.post<Result<Object>>('/volume/save', v)
    }

    prune(node: string) {
        return ajax.post<{
            count: number;
            size: number;
        }>('/volume/prune', { node })
    }
}

export default new VolumeApi
