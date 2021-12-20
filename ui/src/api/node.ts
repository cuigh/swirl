import ajax, { Result } from './ajax'

export interface Node {
    id: string;
    name?: string;
    hostname: string;
    version: number;
    role: string;
    availability: string;
    engineVersion: string;
    arch: string;
    os: string;
    cpu: number;
    memory: number;
    address: string;
    state: string;
    labels?: {
        name: string;
        value: string;
    }[];
    manager?: {
        leader: boolean;
        reachability: string;
        addr: string;
    };
    createdAt: string;
    updatedAt: string;
}

export interface FindResult {
    node: Node;
    raw: string;
}

export class NodeApi {
    find(id: string) {
        return ajax.get<FindResult>('/node/find', { id })
    }

    list(agent: boolean) {
        return ajax.get<Node[]>('/node/list', { agent })
    }

    search() {
        return ajax.get<Node[]>('/node/search')
    }

    save(node: Node) {
        return ajax.post<Result<Object>>('/node/save', node)
    }

    delete(id: string) {
        return ajax.post<Result<Object>>('/node/delete', { id })
    }
}

export default new NodeApi
