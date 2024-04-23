// Copyright (c) 2022 Braden Nicholson

import Settings from "@/views/terminal/settings/Settings.vue";
import Preferences from "@/views/terminal/settings/Preferences.vue";
import Connection from "@/views/terminal/settings/Connection.vue";
import module from "@/views/terminal/settings/module";
import Endpoints from "@/views/terminal/settings/endpoint/Endpoints.vue";
import Zone from "@/views/terminal/settings/zone/Zone.vue";
import Zones from "@/views/terminal/settings/zone/Zones.vue";
import Devices from "@/views/terminal/settings/device/Device.vue";
import DeviceOverview from "@/views/terminal/settings/device/DeviceOverview.vue";
import DeviceMonitor from "@/views/terminal/settings/device/DeviceMonitor.vue";
import RenameDevice from "@/views/terminal/settings/device/RenameDevice.vue";
import subroutines from "@/views/terminal/settings/subroutines";
import entities from "@/views/terminal/settings/entity";
import macros from "@/views/terminal/settings/macro";

export default {
    routes: {
        path: '/terminal/settings',
        name: 'Settings',
        redirect: '/terminal/settings/preferences',
        component: Settings,
        icon: '􀍟',
        meta: {
            icon: "/applications-system.svg"
        },
        children: [
            {
                path: '/terminal/settings/preferences',
                name: 'Preferences',
                icon: '􀐗',
                meta: {
                    order: 0,
                },
                component: Preferences
            },
            {
                path: '/terminal/settings/connection',
                name: 'Connection',
                icon: '􀇃',
                meta: {
                    order: 1,
                },
                component: Connection
            },
            module,
            subroutines,
            {
                path: '/terminal/settings/endpoints',
                name: 'Endpoints',
                icon: '􁅀',
                meta: {
                    order: 4,
                },
                component: Endpoints
            },

            entities,
            {
                path: '/terminal/settings/devices',
                name: 'Devices',
                icon: '􁆬',
                meta: {
                    order: 5,
                },
                component: Devices,
                children: [
                    {
                        path: '/terminal/settings/devices',
                        name: 'Device Overview',
                        icon: 'share-nodes',
                        component: DeviceOverview,
                    },
                    {
                        path: '/terminal/settings/devices/:device/configure',
                        name: 'Edit Device',
                        icon: 'share-nodes',
                        component: RenameDevice,
                    },
                    {
                        path: '/terminal/settings/devices/:device/monitor',
                        name: 'Monitor Device',
                        icon: 'share-nodes',
                        component: DeviceMonitor,
                    },
                ]
            },
            macros,
            {
                path: '/terminal/settings/zones',
                name: 'Zone',
                icon: '􀟻',
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