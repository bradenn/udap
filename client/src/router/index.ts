// Copyright (c) 2022 Braden Nicholson

import {createRouter, createWebHashHistory} from "vue-router";

// Settings views
import settings from "@/views/terminal/settings/index";

import Settings from "@/views/terminal/settings/Settings.vue";
import Preferences from "@/views/terminal/settings/Preferences.vue";
import Connection from "@/views/terminal/settings/Connection.vue";
import Modules from "@/views/terminal/settings/module/Modules.vue";
import Endpoints from "@/views/terminal/settings/endpoint/Endpoints.vue";
import Timings from "@/views/terminal/settings/Timings.vue";
import Zones from "@/views/terminal/settings/zone/Zones.vue";


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
import Whiteboard from "@/views/terminal/whiteboard/Whiteboard.vue";

// Calculator Routes
import Calculator from "@/views/terminal/calculator/Calculator.vue";

// Atlas Routes
import Atlas from "@/views/terminal/atlas/Atlas.vue";
import AtlasSettings from "@/views/terminal/atlas/AtlasSettings.vue";
import AtlasOverview from "@/views/terminal/atlas/AtlasOverview.vue";

// Defense Routes
import Defense from "@/views/terminal/defense/Defense.vue";
import DefenseOverview from "@/views/terminal/defense/DefenseOverview.vue";

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
            icon: 'bars-progress',
            component: Preferences
        },
        {
            path: '/terminal/settings/connection',
            name: 'Connection',
            icon: 'cloud',
            component: Connection
        },
        {
            path: '/terminal/settings/modules',
            name: 'Modules',
            icon: 'layer-group',
            component: Modules
        },
        {
            path: '/terminal/settings/endpoints',
            name: 'Endpoints',
            icon: 'expand',
            component: Endpoints
        },
        {
            path: '/terminal/settings/timings',
            name: 'Timings',
            icon: 'clock',
            component: Timings
        },

        {
            path: '/terminal/settings/zones',
            name: 'Zones',
            icon: 'map',
            component: Zones
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
            icon: 'earth-americas',
            component: Earth,
        },
        {
            path: '/terminal/exogeology/moon',
            name: 'Moon',
            icon: 'moon',
            component: Moon,
        },
        {
            path: '/terminal/exogeology/sol',
            name: 'Sol',
            icon: 'sun',
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
            icon: 'stopwatch',
            component: Stopwatch,
        },
        {
            path: '/terminal/timing/timer',
            name: 'Timer',
            icon: 'clock',
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

const atlasRoutes = {
    path: '/terminal/atlas',
    name: 'Atlas',
    redirect: '/terminal/atlas/overview',
    component: Atlas,
    icon: 'fa-atom',
    children: [
        {
            path: '/terminal/atlas/overview',
            name: 'AtlasOverview',
            component: AtlasOverview,
        },
        {
            path: '/terminal/atlas/settings',
            name: 'AtlasSettings',
            component: AtlasSettings,
        },
    ]
}

const defenseRoutes = {
    path: '/terminal/defense',
    name: 'Sentry',
    redirect: '/terminal/defense/overview',
    component: Defense,
    icon: 'fa-shield',
    children: [
        {
            path: '/terminal/defense/overview',
            name: 'DefenseOverview',
            component: DefenseOverview,
        },
    ]
}

const whiteboardRoutes = {
    path: '/terminal/whiteboard',
    name: 'Whiteboard',
    component: Whiteboard,
    icon: 'fa-highlighter',
}

const calculatorRoute = {
    path: '/terminal/calculator',
    name: 'Calculator',
    component: Calculator,
    icon: 'fa-calculator',
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
        settings.routes,
        timingRoutes,
        exogeologyRoutes,
        weatherRoutes,
        whiteboardRoutes,
        calculatorRoute,
        atlasRoutes,
        defenseRoutes
    ],
}

const routes = [
    defaultRoute,
    setupRoutes,
    terminalRoutes,
    {path: '/:pathMatch(.*)*', redirect: "/"}
]

const router = createRouter({
    history: createWebHashHistory(),
    routes: routes,
})

export default router
