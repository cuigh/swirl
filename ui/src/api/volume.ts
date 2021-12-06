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

export interface PruneResult {
    deletedVolumes: string[];
    reclaimedSpace: number;
}

export class VolumeApi {
    find(name: string) {
        return ajax.get<FindResult>('/volume/find', { name })
    }

    search(args: SearchArgs) {
        return ajax.get<SearchResult>('/volume/search', args)
    }

    delete(name: string) {
        return ajax.post<Result<Object>>('/volume/delete', { name })
    }

    save(v: Volume) {
        return ajax.post<Result<Object>>('/volume/save', v)
    }

    prune() {
        return ajax.post<PruneResult>('/volume/prune')
    }
}

export default new VolumeApi
