import ajax, { Result } from './ajax'

export interface Secret {
    id: string;
    name: string;
    data: string;
    version: number;
    labels?: {
        name: string;
        value: string;
    }[];
    driver: {
        name: string;
        options?: {
            name: string;
            value: string;
        }[];
    }
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
    items: Secret[];
    total: number;
}

export interface FindResult {
    secret: Secret;
    raw: string;
}

export class SecretApi {
    find(id: string) {
        return ajax.get<FindResult>('/secret/find', { id })
    }

    search(args: SearchArgs) {
        return ajax.get<SearchResult>('/secret/search', args)
    }

    save(c: Secret) {
        return ajax.post<Result<Object>>('/secret/save', c)
    }

    delete(id: string) {
        return ajax.post<Result<Object>>('/secret/delete', { id })
    }
}

export default new SecretApi
