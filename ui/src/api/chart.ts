import ajax, { Result } from './ajax'

export interface Chart {
    id: string;
    title: string;
    desc: string;
    metrics: {
        legend: string;
        query: string;
    }[];
    kind?: string;
    dashboard: string;
    type: string;
    unit: string;
    width: number;
    height: number;
    margin: {
        left: number;
        right: number;
        top: number;
        bottom: number;
    };
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

export interface SearchArgs {
    name?: string;
    dashboard?: string;
    pageIndex: number;
    pageSize: number;
}

export interface SearchResult {
    items: Chart[];
    total: number;
}

export class ChartApi {
    search(args: SearchArgs) {
        return ajax.get<SearchResult>('/chart/search', args)
    }

    save(chart: Chart) {
        return ajax.post<Result<Object>>('/chart/save', chart)
    }

    find(id: string) {
        return ajax.get<Chart>('/chart/find', { id })
    }

    delete(id: string, title: string) {
        return ajax.post<Result<Object>>('/chart/delete', { id, title })
    }
}

export default new ChartApi
