import ajax, { Result } from './ajax'

export interface Setting {
    region: RegionSetting;
    ldap: LdapSetting;
    metric: MetricSetting;
}

export interface RegionSetting {
    lang: string;
    timezone: number;
}

export interface LdapSetting {
    enabled: boolean;
    address: string;
    security: number;
    auth: string;
    bind_dn: string;
    bind_pwd: string;
    base_dn: string;
    user_dn: string;
    user_filter: string;
    name_attr: string;
    email_attr: string
}

export interface MetricSetting {
    prometheus: string;
}

export class SettingApi {
    load() {
        return ajax.get<Setting>('/setting/load')
    }

    save(id: string, options: Object) {
        console.log({ id, options })
        return ajax.post<Result<Object>>('/setting/save', { id, options })
    }
}

export default new SettingApi
