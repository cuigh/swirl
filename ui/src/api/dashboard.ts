import ajax, { Result } from './ajax'

export interface Dashboard {
    name: string;
    key: string;
    period: number;
    interval: number;
    charts: ChartInfo[];
}

export interface ChartInfo {
    id: string;
    title: string;
    type: 'line' | 'bar' | 'pie' | 'gauge';
    unit: string;
    width: number;
    height: number;
    margin: {
        left?: number;
        right?: number;
        top?: number;
        bottom?: number;
    };
}

export class DashboardApi {
    fetchData(key: string, charts: string[], period: number) {
        return ajax.get<any>('/dashboard/fetch-data', { key, charts: charts.join(","), period })
    }

    find(name: string, key: string) {
        return ajax.get<Dashboard>('/dashboard/find', { name, key })
    }

    save(dashboard: Dashboard) {
        return ajax.post<Result<Object>>('/dashboard/save', dashboard)
    }
}

export default new DashboardApi
