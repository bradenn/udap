export default {
    path: '/settings',
    name: 'settings',
    component: () => import('./Settings.vue'),
    children: [
        {
            path: '/settings',
            name: 'menu',
            meta: {
                name: "Menu",
                icon: "􀦳",

            },
            component: () => import('@/views/settings/SettingsMenu.vue'),
        },
        {
            path: '/settings/general',
            name: 'general',
            meta: {
                name: "General",
                icon: "􀣋",
                category: "App"
            },
            component: () => import('@/views/settings/General.vue'),
        },
        {
            path: '/settings/connection',
            name: 'connection',
            meta: {
                name: "Connection",
                icon: "􀉣",
                category: "App"
            },
            component: () => import('@/views/settings/Connection.vue'),
        },
        {
            path: '/settings/background',
            name: 'background',
            meta: {
                name: "Background",
                icon: "􀏅",
                category: "App"
            },
            component: () => import('@/views/settings/Background.vue'),
        },
        {
            path: '/settings/config',
            name: 'configuration',
            meta: {
                name: "Configuration",
                icon: "􁇵",
                category: "App"
            },
            component: () => import('@/views/settings/Config.vue'),
        },
        {
            path: '/settings/modules',
            name: 'modules',
            meta: {
                name: "Modules",
                icon: "􁚊",
                category: "System"
            },
            component: () => import('@/views/settings/Modules.vue'),
            children: []
        },
        {
            path: '/settings/modules/:moduleId',
            name: 'module',
            meta: {
                name: "Module",
                icon: "+",
            },
            component: () => import('@/views/settings/module/Module.vue'),
        },
        {
            path: '/settings/zones',
            name: 'zones',
            meta: {
                name: "Zones",
                icon: "􀟼",
                category: "System"
            },
            component: () => import('@/views/settings/Zones.vue'),
        },
        {
            path: '/settings/attributes',
            name: 'attributes',
            meta: {
                name: "Attributes",
                icon: "􀯲",
                category: "System"
            },
            component: () => import('@/views/settings/Attributes.vue'),
        },
        {
            path: '/settings/zones/:zoneId',
            name: 'zone',
            component: () => import('@/views/settings/zones/Zone.vue'),
        },
        {
            path: '/settings/timings',
            name: 'timings',
            meta: {
                name: "Timings",
                icon: "􀐬",
                category: "System"
            },
            component: () => import('@/views/settings/Timings.vue'),
        },
        {
            path: '/settings/webservices',
            name: 'webservices',
            meta: {
                name: "WebServices",
                icon: "􀉣",
                category: "Network"
            },
            component: () => import('@/views/settings/WebServices.vue'),
        },
        {
            path: '/settings/endpoints',
            name: 'endpoints',
            meta: {
                name: "Endpoints",
                icon: "􀏭",
                category: "System"
            },
            component: () => import('@/views/settings/Endpoints.vue'),
        },
        {
            path: '/settings/entities',
            name: 'entity_frame',
            meta: {
                name: "Entities",
                icon: "􁓽",
                category: "System"
            },
            component: () => import('@/views/settings/entities/Frame.vue'),
            children: [
                {
                    path: '/settings/entities',
                    name: 'Entities',
                    component: () => import('@/views/settings/entities/Entities.vue'),
                },
                {
                    path: '/settings/entities/:entityId',
                    name: 'entities_view',
                    component: () => import('@/views/settings/entities/EntityManage.vue'),
                }
            ]
        }
    ]
}
