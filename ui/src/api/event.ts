import ajax, { Result } from './ajax'

export interface Event {
    id: string;
    type: string;
    action: string;
    code: number;
    name: string;
    userId: string;
    username: string;
    time: string;
}

export interface SearchArgs {
    type?: string;
    name?: string;
    pageIndex: number;
    pageSize: number;
}

export interface SearchResult {
    items: Event[];
    total: number;
}

export class EventApi {
    search(args: SearchArgs) {
        return ajax.get<SearchResult>('/event/search', args)
    }
}

export default new EventApi
