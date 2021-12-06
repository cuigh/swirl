import ajax, { Result } from './ajax'

export interface Config {
    id: string;
    name: string;
    data: string;
    version: number;
    labels?: {
        name: string;
        value: string;
    }[];
    templating: {
        name: string;
        options?: {
            name: string;
            value: string;
        }[];
    }
    createdAt: string;
    updatedAt: string;
}

export interface SearchArgs {
    name?: string;
    pageIndex: number;
    pageSize: number;
}

export interface SearchResult {
    items: Config[];
    total: number;
}

export interface FindResult {
    config: Config;
    raw: string;
}

export class ConfigApi {
    find(id: string) {
        return ajax.get<FindResult>('/config/find', { id })
    }

    search(args: SearchArgs) {
        return ajax.get<SearchResult>('/config/search', args)
    }

    save(c: Config) {
        return ajax.post<Result<Object>>('/config/save', c)
    }

    delete(id: string) {
        return ajax.post<Result<Object>>('/config/delete', { id })
    }
}

export default new ConfigApi
