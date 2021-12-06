import ajax, { Result } from './ajax'

export interface Registry {
    id: string;
    name: string;
    url: string;
    username: string;
    password: string;
    createdAt: number;
    updatedAt: number;
}

export class RegistryApi {
    find(id: string) {
        return ajax.get<Registry>('/registry/find', { id })
    }

    search() {
        return ajax.get<Registry[]>('/registry/search')
    }

    save(registry: Registry) {
        return ajax.post<Result<Object>>('/registry/save', registry)
    }

    delete(id: string) {
        return ajax.post<Result<Object>>('/registry/delete',  { id })
    }
}

export default new RegistryApi
