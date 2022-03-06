import {createRouter, createWebHistory} from "vue-router";
import Home from "./views/Home.vue"
import Settings from "./views/Settings.vue";
import Endpoint from "./views/settings/Endpoint.vue";
import Connection from "./views/settings/Connection.vue";
import Entities from "./views/settings/Entities.vue";
import Whiteboard from "./components/apps/Whiteboard.vue";
import Network from "./views/Network.vue";
import Media from "./components/apps/Media.vue";
import Wifi from "./components/apps/Wifi.vue";
import Room from "./views/Room.vue";
import Endpoints from "./views/settings/Endpoints.vue";
import Energy from "./components/apps/Energy.vue";

// Setup view components
import Setup from "./setup/Setup.vue";
import Welcome from "./setup/Controller.vue";
import Authentication from "./setup/Authentication.vue";
import Main from "./views/Main.vue";

const routes = [
    {
        path: '/',
        redirect: '/setup/controller',
    },
    {
        path: '/setup',
        name: 'Setup',
        redirect: '/setup/controller',
        component: Setup,
        children: [
            {
                path: '/setup/authentication',
                name: 'Authenticate',
                component: Authentication
            },
            {
                path: '/setup/controller',
                name: 'Welcome',
                component: Welcome
            },
        ]
    },
    {
        path: '/terminal',
        name: 'Terminal',
        redirect: '/terminal/home',
        component: Main,
        children: [
            {
                path: '/terminal/apps/whiteboard',
                name: 'Whiteboard',
                component: Whiteboard,
                icon: '􀏒',
                meta: {slideOrder: 4},
            },
            {
                path: '/terminal/network',
                name: 'Network',
                component: Network,
                icon: '􀏒',
                meta: {slideOrder: 4},
            },
            {
                path: '/terminal/wifi',
                name: 'Wifi',
                component: Wifi,
                icon: '􀏒',
                meta: {slideOrder: 4},
            },
            {
                path: '/terminal/apps/media',
                name: 'Media',
                component: Media,
                icon: '􀑪',
                meta: {slideOrder: 4},
            },
            {
                path: '/terminal/apps/room',
                name: 'Room',
                component: Room,
                icon: '􀑪',
                meta: {slideOrder: 4},
            },
            {
                path: '/terminal/energy',
                name: 'Energy',
                component: Energy,
                icon: '',
                meta: {slideOrder: 4},
            },
            {
                path: '/terminal/home',
                name: 'Home',
                component: Home,
                meta: {slideOrder: 0}
            },
            {
                path: '/terminal/settings',
                name: 'Settings',
                redirect: '/terminal/settings/endpoint',
                component: Settings,
                meta: {slideOrder: 4},
                children: [
                    {
                        path: '/terminal/settings/endpoint',
                        name: 'Endpoint',
                        component: Endpoint,
                        icon: 'bi-display',
                        sf: '􀢹',
                        meta: {slideOrder: 4},
                    },
                    {
                        path: '/terminal/settings/endpoints',
                        name: 'Endpoints',
                        component: Endpoints,
                        icon: 'bi-display',
                        sf: '􀞿',
                        meta: {slideOrder: 4},
                    },
                    {
                        path: '/terminal/settings/entities',
                        name: 'Entities',
                        icon: 'bi-lightbulb',
                        sf: '􀛭',
                        component: Entities,
                        meta: {slideOrder: 4},
                    },
                    {
                        path: '/terminal/settings/connection',
                        name: 'Connection',
                        icon: 'bi-link',
                        sf: '􀉣',
                        component: Connection,
                        meta: {slideOrder: 4},
                    },
                ]
            },

        ]
    },

]


const router = createRouter({history: createWebHistory(), routes})

export default router
