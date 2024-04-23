// Copyright (c) 2023 Braden Nicholson

export interface AppMeta {
    name: string
    icon: string
    link: string
    child?: boolean
    section?: string

}

export {}

export default {
    path: '/apps',
    name: 'apps',
    component: () => import('./Apps.vue'),
    children: [
        {
            path: '/apps',
            name: 'appsMain',
            component: () => import("@/views/apps/AppsView.vue"),
            meta: {
                name: "Apps",
                icon: "􀙣"
            }
        },
        {
            path: '/apps/request',
            name: 'request',
            component: () => import("@/views/apps/RequestBuilder.vue"),
            meta: {
                name: "Request Builder",
                icon: "􀙣"
            }
        },
        {
            path: '/apps/message',
            name: 'message',
            component: () => import("@/views/apps/Message.vue"),
            meta: {
                child: true,
                name: "Messages",
                icon: "􀈠"
            }
        },
        {
            path: '/apps/todo',
            name: 'request',
            component: () => import("@/views/apps/Todo.vue"),
            meta: {
                child: true,
                name: "Todo",
                icon: "􀷾"
            }
        },
        {
            path: '/apps/sensors',
            name: 'sensors',
            component: () => import("@/views/apps/Sensors.vue"),
            meta: {
                name: "Sensors",
                child: true,
                link: "/apps/sensors",
                icon: '􁔊'
            }
        },
        {
            path: '/apps/uhidv2',
            name: 'uhidv2',
            component: () => import("@/views/apps/UHIDv2App.vue"),
            meta: {
                name: "UHID",
                link: "/apps/uhidv2",
                section: "environment",
                icon: '􁔊'
            }
        },
        {
            path: '/apps/uhidv2/:entityId',
            name: 'uhidv2-device',
            component: () => import("@/views/apps/UHIDv2Device.vue"),

        },
        {
            path: '/apps/actions',
            name: 'actions',
            component: () => import("@/views/apps/Actions.vue"),
            icon: '􀎟',
            meta: {
                name: "Actions",
                link: "/apps/actions",
                section: "system",
                icon: '􀒘'
            },
        },
        {
            path: '/apps/environment',
            name: 'actions',
            component: () => import("@/views/apps/Environment.vue"),
            icon: '􀎟',
            meta: {
                name: "Environment",
                link: "/apps/environment",
                section: "environment",
                icon: '􀒘'
            },
        },
        {
            path: '/apps/conversions',
            name: 'conversions',
            component: () => import("@/views/apps/Conversions.vue"),
            meta: {
                name: "Convert",
                link: "/apps/conversions",
                section: "tools",
                icon: '􀅮'
            },
        },
        {
            path: '/apps/library',
            name: 'library',
            component: () => import("@/views/apps/library/Library.vue"),
            meta: {
                name: "Library",
                link: "/apps/Library",
                section: "media",
                icon: '􀬓'
            },
        },
        {
            path: '/apps/actions/create',
            name: 'actionsCreate',
            component: () => import("@/views/apps/ActionCreate.vue"),
            icon: '􀎟',
        },
        {
            path: '/apps/actions/:actionId',
            name: 'action',
            component: () => import("@/views/apps/Action.vue"),
            icon: '􀎟',
        },
        {
            path: '/apps/sensors/:sensorId',
            name: 'sensorView',
            component: () => import("@/views/apps/Sensor.vue"),
            icon: '􀎟',
        },
        {
            path: '/apps/beam',
            name: 'beam',
            component: () => import("@/views/beam/Beam.vue"),
            icon: '􀎟',

        },
        {
            path: '/apps/movies',
            name: 'movies',
            component: () => import("@/views/apps/Movies.vue"),
            meta: {
                name: "Movies",
                link: "/apps/movies",
                section: "media",
                icon: '􀜤'
            },
        },
        {
            path: '/apps/network',
            name: 'network',
            component: () => import("@/views/apps/network/Network.vue"),
            meta: {
                name: "Network",
                link: "/apps/network",
                section: "system",
                icon: '􁆬'
            },
        },
        {
            path: '/apps/network/subdomains',
            name: 'subdomains',
            component: () => import("@/views/apps/network/Subdomains.vue"),
            icon: '􀎟',
        },
        {
            path: '/apps/network/graph',
            name: 'graphs',
            component: () => import("@/views/apps/network/Graph.vue"),
            icon: '􀎟',
        },
        {
            path: '/apps/sentry',
            name: 'sentry',
            component: () => import("@/views/beam/Beam.vue"),
            meta: {
                name: "Lasers",
                link: "/sentry",
                section: "environment",
                icon: '􀑃'
            },
        },
        {
            path: '/apps/monitor',
            name: 'monitor',
            component: () => import("@/views/monitor/MonitorList.vue"),
            icon: '􀎟',
            meta: {
                name: "Monitor",
                link: "/monitor",
                section: "system",
                icon: '􁂥'
            },
        },
        {
            path: '/apps/lights',
            name: 'lights',
            component: () => import("@/views/apps/Lights.vue"),
            meta: {
                name: "Lights",
                link: "/lights",
                section: "environment",
                icon: '􁓼'
            },
        },
        {
            path: '/apps/thermostat',
            name: 'thermostat',
            component: () => import("@/views/home/ThermostatView.vue"),
            meta: {
                name: "Thermostat",
                link: "/thermostat",
                section: "environment",
                icon: '􁁋'
            },
        },
        {
            path: '/apps/detour',
            name: 'detour',
            component: () => import("@/views/apps/Detour.vue"),
            meta: {
                name: "Detour",
                link: "/apps/detour",
                section: "tools",
                icon: '􁎱'
            },
        },
        {
            path: '/apps/utilities',
            name: 'utilities',
            component: () => import("@/views/apps/utilities/Utilities.vue"),
            children: [
                {
                    path: '/apps/utilities',
                    name: 'utilitiesMenu',
                    component: () => import("@/views/apps/utilities/UtilitiesMenu.vue"),
                    meta: {
                        name: "Utilities",
                        icon: "􀙣",
                        child: true
                    }
                },
                {
                    path: '/apps/utilities/torque',
                    name: 'utilitiesTorque',
                    component: () => import("@/views/apps/utilities/Torque.vue"),
                    meta: {
                        name: "Torque",
                        icon: "􀙣",
                        child: true
                    }
                },
            ],
            meta: {
                name: "Utilities",
                link: "/utilities",
                section: "tools",
                icon: '􀤋'
            },
        },
        {
            path: '/apps/balance',
            name: 'balanceView',
            component: () => import("@/views/apps/budget/BudgetView.vue"),
            children: [
                {
                    path: '/apps/balance',
                    name: 'balance',
                    component: () => import("@/views/apps/budget/Balance.vue"),
                    meta: {
                        name: "Expense",
                        icon: "􀙣",
                        child: true
                    }
                },
                {
                    path: '/apps/balance/expense',
                    name: 'buildExpense',
                    component: () => import("@/views/apps/budget/BuildExpense.vue"),
                    meta: {
                        name: "Build Expense",
                        icon: "􀙣",
                        child: true
                    }
                },
            ],
            meta: {
                name: "Balance",
                link: "/balance",
                section: "tools",
                icon: '􁎣'
            },
        },
    ]
}