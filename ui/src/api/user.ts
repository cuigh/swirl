import ajax, { Result } from './ajax'

export interface AuthUser {
    token: string;
    id: string;
    name: string;
}

export interface User {
    id: string;
    name: string;
    loginName: string;
    password: string;
    passwordConfirm?: string;
    admin: boolean;
    type: string;
    status: number;
    email: string;
    roles: string[];
    createdAt: number;
    updatedAt: number;
    createdBy: {
        id: string;
        name: string;
    };
    updatedBy: {
        id: string;
        name: string;
    };
}

export interface LoginArgs {
    name: string;
    password: string;
}

export interface SearchArgs {
    name?: string;
    loginName?: string;
    filter?: string;
    pageIndex: number;
    pageSize: number;
}

export interface SearchResult {
    items: User[];
    total: number;
}

export interface SetStatusArgs {
    id: string;
    status: number;
}

export interface ModifyPasswordArgs {
    oldPwd: string;
    newPwd: string;
}

export class UserApi {
    login(args: LoginArgs) {
        return ajax.post<AuthUser>('/user/sign-in', args)
    }

    save(user: User) {
        return ajax.post<Result<Object>>('/user/save', user)
    }

    find(id: string) {
        return ajax.get<User>('/user/find', { id })
    }

    fetch(ids: string[]) {
        return ajax.get<User[]>('/user/fetch', { ids: ids.join(',') })
    }

    search(args: SearchArgs) {
        return ajax.get<SearchResult>('/user/search', args)
    }

    setStatus(args: SetStatusArgs) {
        return ajax.post<Result<Object>>('/user/set-status', args)
    }

    delete(id: string, name: string) {
        return ajax.post<Result<Object>>('/user/delete', { id, name })
    }

    modifyPassword(args: ModifyPasswordArgs) {
        return ajax.post<Result<Object>>('/user/modify-password', args)
    }

    modifyProfile(user: User) {
        return ajax.post<Result<Object>>('/user/modify-profile', user)
    }
}

export default new UserApi