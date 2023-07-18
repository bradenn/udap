// Copyright (c) 2022 Braden Nicholson

import {createRouter, createWebHistory} from "vue-router";

import setupRoutes from "@/views/setup";


import settings from "@/views/settings";
import Demo from "@/views/Demo.vue";


const demoRoute = {
    path: '/home/demo',
    name: 'demo',
    component: Demo,
    icon: '􀎟',
}


const mobileRoutes =
    {
        path: '/home',
        name: 'Base',
        redirect: "/home/dashboard",
        component: () => import("@/views/Dashboard.vue"),
        icon: '􀎟',
        children: [
            {
                path: '/home/dashboard',
                name: 'dashboard',
                component: () => import("@/views/Home.vue"),
                icon: '􀎟',
            },
            {
                path: '/home/apps',
                name: 'apps',
                component: () => import("@/views/apps/Apps.vue"),
                icon: '􀎟',
            }, {
                path: '/home/beampath',
                name: 'beampath',
                component: () => import("@/views/apps/BeamPath.vue"),
                icon: '􀎟',
            },
            {
                path: '/home/entity/:entityId',
                name: 'editEntity',
                component: () => import("@/views/home/EntityEdit.vue"),
                icon: '􀎟',
            },
            {
                path: '/home/extra/homework',
                name: 'homework',
                component: () => import("@/views/extra/Homework.vue"),
                icon: '􀎟',
            },
            {
                path: '/home/subroutines',
                name: 'subroutines',
                component: () => import("@/views/home/SubRoutineList.vue"),
                icon: '􀎟',
            },
            {
                path: '/home/sentry',
                name: 'sentry',
                component: () => import("@/views/beam/Beam.vue"),
                icon: '􀎟',
            },
            {
                path: '/home/monitor',
                name: 'monitor',
                component: () => import("@/views/monitor/MonitorList.vue"),
                icon: '􀎟',
            },
            {
                path: '/home/lights',
                name: 'lights',
                component: () => import("@/views/apps/Lights.vue"),
                icon: '􀎟',
            },
            {
                path: '/home/menu',
                name: 'menu',
                component: () => import("@/views/MenuPage.vue"),
                icon: '􀎟',
            },
            {
                path: '/home/thermostat',
                name: 'thermostat',
                component: () => import("@/views/home/ThermostatView.vue"),
                icon: '􀎟',
            },
            {
                path: '/home/macros',
                name: 'macros',
                component: () => import("@/views/home/MacroList.vue"),
                icon: '􀎟',
            },
            {
                path: '/home/beam',
                name: 'beam',
                component: () => import("@/views/beam/Beam.vue"),
                icon: '􀎟',
            },
            settings,
            demoRoute
        ]
    };

const routes = [
    mobileRoutes,
    setupRoutes,

    // {path: '/:pathMatch(.*)*', redirect: "/"}
]

const router = createRouter({
    history: createWebHistory("/"),
    routes: routes,
})

export default router
