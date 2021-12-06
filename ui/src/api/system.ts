import ajax, { Result } from './ajax'

export interface User {
    loginName: string;
    password: string;
    name: string;
    email: string;
}

export interface State {
    fresh: boolean;
}

export interface Version {
    version: string;
    goVersion: string;
}

export interface Summary {
    version: string;
    goVersion: string;
    nodeCount: number;
    networkCount: number;
    serviceCount: number;
    stackCount: number;
}

export class SystemApi {
    checkState() {
        return ajax.get<State>('/system/check-state')
    }

    createAdmin(user: User) {
        return ajax.post<Result<Object>>('/system/create-admin', user)
    }

    version() {
        return ajax.get<Version>('/system/version')
    }

    summarize() {
        return ajax.get<Summary>('/system/summarize')
    }
}

export default new SystemApi
