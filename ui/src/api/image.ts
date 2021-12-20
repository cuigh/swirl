import ajax, { Result } from './ajax'

export interface Image {
    id: string;
    pid: string;
    created: string;
    containers: number;
    digests: string[];
    tags: string[];
    labels?: {
        name: string;
        value: string;
    }[];
    size: number;
    sharedSize: number;
    virtualSize: number;
    comment: string;
    container: string;
    dockerVersion: string;
    author: string;
    variant: string;
    arch: string;
    os: string;
    osVersion: string;
    lastTagTime: string;
    graphDriver?: {
        name?: string;
        data?: {
            name: string;
            value: string;
        }[];
    };
    rootFS?: {
        type?: string;
        layers?: string[];
        baseLayer?: string;
    };
    histories?: {
        id: string;
        comment: string;
        size: number;
        tags: string[];
        createdAt: string;
        createdBy: string;
    }[];
}

export interface SearchArgs {
    node?: string;
    name?: string;
    pageIndex: number;
    pageSize: number;
}

export interface SearchResult {
    items: Image[];
    total: number;
}

export interface FindResult {
    image: Image;
    raw: string;
}

export class ImageApi {
    find(node: string, id: string) {
        return ajax.get<FindResult>('/image/find', { node, id })
    }

    search(args: SearchArgs) {
        return ajax.get<SearchResult>('/image/search', args)
    }

    delete(node: string, id: string, name: string) {
        return ajax.post<Result<Object>>('/image/delete', { node, id, name })
    }
}

export default new ImageApi
