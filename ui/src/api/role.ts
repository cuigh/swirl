import ajax, { Result } from './ajax'

export interface Role {
    id: string;
    name: string;
    desc: string;
    perms: string[];
    createdAt: string;
    updatedAt: string;
    createdBy: {
        id: string;
        name: string;
    };
    updatedBy: {
        id: string;
        name: string;
    };
}

export class RoleApi {
    find(id: string) {
        return ajax.get<Role>('/role/find', { id })
    }

    search(name?: string) {
        return ajax.get<Role[]>('/role/search', { name })
    }

    save(role: Role) {
        return ajax.post<Result<Object>>('/role/save', role)
    }

    delete(id: string, name: string) {
        return ajax.post<Result<Object>>('/role/delete',  { id, name })
    }
}

export default new RoleApi
