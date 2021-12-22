export const perms = [
    {
        key: 'registry',
        actions: ['view', 'edit', 'delete'],
    },
    {
        key: 'node',
        actions: ['view', 'edit', 'delete'],
    },
    {
        key: 'network',
        actions: ['view', 'edit', 'delete', 'disconnect'],
    },
    {
        key: 'service',
        actions: ['view', 'edit', 'delete', 'restart', 'rollback', 'logs'],
    },
    {
        key: 'task',
        actions: ['view', 'logs'],
    },
    {
        key: 'stack',
        actions: ['view', 'edit', 'delete', 'deploy', 'shutdown'],
    },
    {
        key: 'config',
        actions: ['view', 'edit', 'delete'],
    },
    {
        key: 'secret',
        actions: ['view', 'edit', 'delete'],
    },
    {
        key: 'image',
        actions: ['view', 'delete'],
    },
    {
        key: 'container',
        actions: ['view', 'delete', 'logs', 'execute'],
    },
    {
        key: 'volume',
        actions: ['view', 'edit', 'delete'],
    },
    {
        key: 'user',
        actions: ['view', 'edit', 'delete'],
    },
    {
        key: 'role',
        actions: ['view', 'edit', 'delete'],
    },
    {
        key: 'chart',
        actions: ['view', 'edit', 'delete'],
    },
    {
        key: 'dashboard',
        actions: ['edit'],
    },
    {
        key: 'event',
        actions: ['view'],
    },
    {
        key: 'setting',
        actions: ['view', 'edit'],
    },
]
