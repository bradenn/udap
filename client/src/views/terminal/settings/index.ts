// Copyright (c) 2022 Braden Nicholson

import Settings from "@/views/terminal/settings/Settings.vue";
import Preferences from "@/views/terminal/settings/Preferences.vue";
import Connection from "@/views/terminal/settings/Connection.vue";
import Modules from "@/views/terminal/settings/module/Modules.vue";
import Endpoints from "@/views/terminal/settings/endpoint/Endpoints.vue";
import Timings from "@/views/terminal/settings/Timings.vue";
import Zone from "@/views/terminal/settings/zone/Zone.vue";
import Zones from "@/views/terminal/settings/zone/Zones.vue";
import Entities from "@/views/terminal/settings/entity/Entities.vue";

export default {
    routes: {
        path: '/terminal/settings',
        name: 'Settings',
        redirect: '/terminal/settings/preferences',
        component: Settings,
        icon: 'fa-cog',
        children: [
            {
                path: '/terminal/settings/preferences',
                name: 'Preferences',
                icon: 'bars-progress',

                meta: {
                    order: 0,
                },
                component: Preferences
            },
            {
                path: '/terminal/settings/connection',
                name: 'Connection',
                icon: 'cloud',
                meta: {
                    order: 1,
                },
                component: Connection
            },
            {
                path: '/terminal/settings/modules',
                name: 'Modules',
                icon: 'layer-group',
                meta: {
                    order: 2,
                },
                component: Modules
            },
            {
                path: '/terminal/settings/endpoints',
                name: 'Endpoints',
                icon: 'expand',
                meta: {
                    order: 3,
                },
                component: Endpoints
            },
            {
                path: '/terminal/settings/entities',
                name: 'Entities',
                icon: 'expand',
                meta: {
                    order: 4,
                },
                component: Entities
            },
            {
                path: '/terminal/settings/timings',
                name: 'Timings',
                icon: 'clock',
                meta: {
                    order: 5,
                },
                component: Timings
            },


            {
                path: '/terminal/settings/zones',
                name: 'Zone',
                icon: 'map',
                meta: {
                    order: 6,
                },
                component: Zone,
                children: [
                    {
                        path: '/terminal/settings/zones',
                        name: 'Zones',
                        icon: 'map',
                        component: Zones,
                    },
                ]
            },
        ]
    }
}