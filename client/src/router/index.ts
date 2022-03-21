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

// Exogeology Routes
import Exogeology from "../views/terminal/exogeology/Exogeology.vue";
import Earth from "../views/terminal/exogeology/Earth.vue";
import Moon from "../views/terminal/exogeology/Moon.vue";
import Sol from "../views/terminal/exogeology/Sol.vue";

// Timing Routes
import Timing from "../views/terminal/timing/Timing.vue";
import Stopwatch from "../views/terminal/timing/Stopwatch.vue";
import Timer from "../views/terminal/timing/Timer.vue";

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

const exogeologyRoutes = {
    path: '/terminal/exogeology',
    name: 'Exogeology',
    redirect: '/terminal/exogeology/earth',
    component: Exogeology,
    children: [
        {
            path: '/terminal/exogeology/earth',
            name: 'Earth',
            component: Earth,
        },
        {
            path: '/terminal/exogeology/moon',
            name: 'Moon',
            component: Moon,
        },
        {
            path: '/terminal/exogeology/sol',
            name: 'Sol',
            component: Sol,
        },
    ]
}
const timingRoutes = {
    path: '/terminal/timing',
    name: 'Timing',
    redirect: '/terminal/timing/stopwatch',
    component: Timing,
    children: [
        {
            path: '/terminal/timing/stopwatch',
            name: 'Stopwatch',
            component: Stopwatch,
        },
        {
            path: '/terminal/timing/timer',
            name: 'Timer',
            component: Timer,
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
        settingsRoutes,
        timingRoutes,
        exogeologyRoutes
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
