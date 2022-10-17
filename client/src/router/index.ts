// Copyright (c) 2022 Braden Nicholson

import {createRouter, createWebHashHistory} from "vue-router";

// Settings views
import settings from "@/views/terminal/settings/index";
import setup from "@/views/setup";
import haptics from "@/views/tests/Haptics.vue"

// Setup view components
import Setup from "@/views/setup/Setup.vue";

// Terminal Routes
import Terminal from "../views/terminal/Terminal.vue";
import Home from "../views/terminal/Home.vue"
import Wifi from "../views/terminal/wifi/Wifi.vue";
import Energy from "../views/terminal/Energy.vue";

// Exogeology Routes
import Exogeology from "../views/terminal/exogeology/Exogeology.vue";
import Earth from "../views/terminal/exogeology/Earth.vue";
import Magnetosphere from "../views/terminal/exogeology/Magnetosphere.vue";
import Moon from "../views/terminal/exogeology/Moon.vue";
import Sol from "../views/terminal/exogeology/Sol.vue";

// Timing Routes
import Timing from "../views/terminal/timing/Timing.vue";
import Stopwatch from "../views/terminal/timing/Stopwatch.vue";
import Timer from "../views/terminal/timing/Timer.vue";

// Weather
import weather from "@/views/terminal/weather"
// Atlas Routes
import Atlas from "@/views/terminal/atlas/Atlas.vue";
import AtlasSettings from "@/views/terminal/atlas/AtlasSettings.vue";
import AtlasOverview from "@/views/terminal/atlas/AtlasOverview.vue";

// Defense Routes
import Defense from "@/views/terminal/defense/Defense.vue";
import DefenseOverview from "@/views/terminal/defense/DefenseOverview.vue";

// Layout Routes
import Layout from "@/views/terminal/layout/Layout.vue";

// Layout Routes
import LayoutOverview from "@/views/terminal/layout/LayoutOverview.vue";

import Aliens from "@/views/terminal/aliens/Aliens.vue";
import Drake from "@/views/terminal/aliens/Drake.vue";
import diagnostic from "@/views/terminal/diagnostic";
import periodic from "@/views/terminal/periodic";
import remote from "@/views/terminal/remote";


const Whiteboard = () => import("@/views/terminal/whiteboard/Whiteboard.vue")

// Calculator Routes
const Calculator = () => import("@/views/terminal/calculator/Calculator.vue");

const defaultRoute = {
    path: '/',
    component: Setup,
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
        }, {
            path: '/terminal/exogeology/magnetosphere',
            name: 'Magnetosphere',
            icon: 'earth-americas',
            component: Magnetosphere,
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


const alienRoutes = {
    path: '/terminal/aliens',
    name: 'Drake',
    redirect: '/terminal/aliens/drake',
    component: Aliens,
    icon: 'fa-rocket',
    children: [
        {
            path: '/terminal/aliens/drake',
            name: 'AliensDrake',
            component: Drake,
            icon: 'fa-house',
        },
    ]
};

const layoutRoutes = {
    path: '/terminal/layout',
    name: 'Layout',
    component: Layout,
    redirect: '/terminal/layout/overview',
    icon: 'fa-calculator',
    children: [
        {
            path: '/terminal/layout/overview',
            name: 'LayoutOverview',
            component: LayoutOverview,
            icon: 'fa-house',
        },
    ]
};

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
            name: 'Sentry Overview',
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
    meta: {
        status: ''
    },
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
        diagnostic,
        weather,
        remote,
        whiteboardRoutes,
        calculatorRoute,
        atlasRoutes,
        alienRoutes,
        defenseRoutes,
        layoutRoutes,
        periodic
    ],
}


const testRoutes = {
    path: '/test',
    name: 'TESTRoute',
    component: haptics,
}

const routes = [
    defaultRoute,
    setup,
    terminalRoutes,
    {path: '/:pathMatch(.*)*', redirect: "/"}
]

const router = createRouter({
    history: createWebHashHistory(),
    routes: routes,
})

export default router
