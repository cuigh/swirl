import ajax, { Result } from './ajax'

export interface Stack {
    id: string;
    name: string;
    content: string;
    services?: string[];
    internal: boolean;
    createdAt: string;
    createdBy: string;
    updatedAt: string;
    updatedBy: string;
}

export interface SearchArgs {
    name?: string;
    filter?: string;
}

export class StackApi {
    find(name: string) {
        return ajax.get<Stack>('/stack/find', { name })
    }

    search(args: SearchArgs) {
        return ajax.get<Stack[]>('/stack/search', args)
    }

    save(s: Stack) {
        return ajax.post<Result<Object>>('/stack/save', s)
    }

    delete(name: string) {
        return ajax.post<Result<Object>>('/stack/delete', { name })
    }

    shutdown(name: string) {
        return ajax.post<Result<Object>>('/stack/shutdown', { name })
    }

    deploy(name: string) {
        return ajax.post<Result<Object>>('/stack/deploy', { name })
    }
}

export default new StackApi
