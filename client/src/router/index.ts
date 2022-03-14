// Copyright (c) 2022 Braden Nicholson

import {createRouter, createWebHistory} from "vue-router";

// Settings views
import Settings from "@/views/settings/Settings.vue";
import Preferences from "@/views/settings/Preferences.vue";
import Connection from "@/views/settings/Connection.vue";

// Setup view components
import Setup from "@/views/setup/Setup.vue";
import Welcome from "@/views/setup/Controller.vue";
import Authentication from "@/views/setup/Authentication.vue";

// Terminal Routes
import Terminal from "../views/terminal/Terminal.vue";
import Home from "../views/terminal/Home.vue"
import Wifi from "../views/terminal/Wifi.vue";
import Energy from "../views/terminal/Energy.vue";

const defaultRoute = {
    path: '/',
    component: Setup,
}

const setupRoutes = {
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
}

const settingsRoutes = {
    path: '/terminal/settings',
    name: 'Settings',
    redirect: '/terminal/settings/preferences',
    component: Settings,
    children: [
        {
            path: '/terminal/settings/preferences',
            name: 'Preferences',
            component: Preferences
        },
        {
            path: '/terminal/settings/connection',
            name: 'Connection',
            component: Connection
        },
    ]
}

const terminalRoutes = {
    path: '/terminal',
    name: 'Terminal',
    redirect: '/terminal/home',
    component: Terminal,
    children: [
        {
            path: '/terminal/home',
            name: 'Home',
            component: Home,
        },
        {
            path: '/terminal/wifi',
            name: 'Wifi',
            component: Wifi,
        },
        {
            path: '/terminal/energy',
            name: 'Energy',
            component: Energy,
        },
        settingsRoutes
    ]
}

const routes = [
    defaultRoute,
    setupRoutes,
    terminalRoutes
]

const router = createRouter({
    history: createWebHistory(),
    routes: routes,
})

export default router
