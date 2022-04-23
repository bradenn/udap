// Copyright (c) 2022 Braden Nicholson

import {createRouter, createWebHistory} from "vue-router";

// Settings views
import Settings from "@/views/terminal/settings/Settings.vue";
import Preferences from "@/views/terminal/settings/Preferences.vue";
import Connection from "@/views/terminal/settings/Connection.vue";

// Setup view components
import Setup from "@/views/setup/Setup.vue";
import Welcome from "@/views/setup/Controller.vue";
import Authentication from "@/views/setup/Authentication.vue";

// Terminal Routes
import Terminal from "../views/terminal/Terminal.vue";
import Home from "../views/terminal/Home.vue"
import Wifi from "../views/terminal/wifi/Wifi.vue";
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

// Weather Routes
import WeatherApp from "../views/terminal/weather/Weather.vue";
import Summary from "../views/terminal/weather/Summary.vue";


// Whiteboard Routes
import Whiteboard from "../views/terminal/whiteboard/Whiteboard.vue";

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
    icon: 'fa-cog',
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
    name: 'ExoGeo',
    redirect: '/terminal/exogeology/earth',
    component: Exogeology,
    icon: 'fa-satellite',
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
    icon: 'fa-stopwatch',
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

const weatherRoutes = {
    path: '/terminal/weather',
    name: 'Weather',
    redirect: '/terminal/weather/summary',
    component: WeatherApp,
    icon: 'fa-cloud-sun',
    children: [
        {
            path: '/terminal/weather/summary',
            name: 'Summary',
            component: Summary,
        },
    ]
}

const whiteboardRoutes = {
    path: '/terminal/whiteboard',
    name: 'Whiteboard',
    component: Whiteboard,
    icon: 'fa-highlighter',
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
            icon: 'fa-house',
        },
        {
            path: '/terminal/wifi',
            name: 'Wifi',
            component: Wifi,
            icon: 'fa-wifi',
        },
        {
            path: '/terminal/energy',
            name: 'Energy',
            component: Energy,
            icon: 'fa-bolt',
        },
        settingsRoutes,
        timingRoutes,
        exogeologyRoutes,
        weatherRoutes,
        whiteboardRoutes
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
