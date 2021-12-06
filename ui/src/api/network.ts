import ajax, { Result } from './ajax'

export interface Network {
    id: string;
    name: string;
    created: string;
    scope: string;
    driver: string;
    internal: boolean;
    attachable: boolean;
    ingress: boolean;
    ipv6: boolean;
    ipam: {
        driver: string;
        options: [];
        config: {
            subnet: string;
            gateway: string;
            range: string;
        }[],
    }
    options?: {
        name: string;
        value: string;
    }[];
    labels?: {
        name: string;
        value: string;
    }[];
    containers?: {
        id: string;
        name: string;
        mac: string;
        ipv4: string;
        ipv6: string;
    }[];
}

export interface FindResult {
    network: Network;
    raw: string;
}

export class NetworkApi {
    find(name: string) {
        return ajax.get<FindResult>('/network/find', { name })
    }

    search() {
        return ajax.get<Network[]>('/network/search')
    }

    save(network: Network) {
        return ajax.post<Result<Object>>('/network/save', network)
    }

    delete(id: string, name: string) {
        return ajax.post<Result<Object>>('/network/delete', { id, name })
    }

    disconnect(networkId: string, networkName: string, container: string) {
        return ajax.post<Result<Object>>('/network/disconnect', { networkId, networkName, container })
    }
}

export default new NetworkApi
