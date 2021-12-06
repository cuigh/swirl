export const perms = [
    {
        key: 'registry',
        items: [
            { key: "registry.view", perm: "view" },
            { key: "registry.edit", perm: "edit" },
            { key: "registry.delete", perm: "delete" },
        ],
    },
    {
        key: 'node',
        items: [
            { key: "node.view", perm: "view" },
            { key: "node.edit", perm: "edit" },
            { key: "node.delete", perm: "delete" },
        ],
    },
    {
        key: 'network',
        items: [
            { key: "network.view", perm: "view" },
            { key: "network.edit", perm: "edit" },
            { key: "network.delete", perm: "delete" },
            { key: "network.disconnect", perm: "disconnect" },
        ],
    },
    {
        key: 'service',
        items: [
            { key: "service.view", perm: "view" },
            { key: "service.edit", perm: "edit" },
            { key: "service.delete", perm: "delete" },
            { key: "service.restart", perm: "restart" },
            { key: "service.rollback", perm: "rollback" },
            { key: "service.logs", perm: "logs" },
        ],
    },
    {
        key: 'task',
        items: [
            { key: "task.view", perm: "view" },
            { key: "task.logs", perm: "logs" },
        ],
    },
    {
        key: 'stack',
        items: [
            { key: "stack.view", perm: "view" },
            { key: "stack.edit", perm: "edit" },
            { key: "stack.delete", perm: "delete" },
            { key: "stack.deploy", perm: "deploy" },
            { key: "stack.shutdown", perm: "shutdown" },
        ],
    },
    {
        key: 'config',
        items: [
            { key: "config.view", perm: "view" },
            { key: "config.edit", perm: "edit" },
            { key: "config.delete", perm: "delete" },
        ],
    },
    {
        key: 'secret',
        items: [
            { key: "secret.view", perm: "view" },
            { key: "secret.edit", perm: "edit" },
            { key: "secret.delete", perm: "delete" },
        ],
    },
    {
        key: 'image',
        items: [
            { key: "image.view", perm: "view" },
            { key: "image.delete", perm: "delete" },
        ],
    },
    {
        key: 'container',
        items: [
            { key: "container.view", perm: "view" },
            { key: "container.delete", perm: "delete" },
            { key: "container.logs", perm: "logs" },
        ],
    },
    {
        key: 'volume',
        items: [
            { key: "volume.view", perm: "view" },
            { key: "volume.edit", perm: "edit" },
            { key: "volume.delete", perm: "delete" },
        ],
    },
    {
        key: 'user',
        items: [
            { key: "user.view", perm: "view" },
            { key: "user.edit", perm: "edit" },
            { key: "user.delete", perm: "delete" },
        ],
    },
    {
        key: 'role',
        items: [
            { key: "role.view", perm: "view" },
            { key: "role.edit", perm: "edit" },
            { key: "role.delete", perm: "delete" },
        ],
    },
    {
        key: 'chart',
        items: [
            { key: "chart.view", perm: "view" },
            { key: "chart.edit", perm: "edit" },
            { key: "chart.delete", perm: "delete" },
            { key: "chart.dashboard", perm: "dashboard" },
        ],
    },
    {
        key: 'event',
        items: [
            { key: "event.view", perm: "view" },
        ],
    },
    {
        key: 'setting',
        items: [
            { key: "setting.view", perm: "view" },
            { key: "setting.edit", perm: "edit" },
        ],
    },
]

function contains(arr1: string[], arr2: string[]): boolean {
    const set = new Set(arr1);
    return arr2.every(s => set.has(s))
}