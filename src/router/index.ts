// Copyright (c) 2022 Braden Nicholson

import {createRouter, createWebHistory} from "vue-router";

import setupRoutes from "@/views/setup";


import settings from "@/views/settings";
import NotFound from "@/views/NotFound.vue";


const mobileRoutes =
    {
        path: '/',
        name: 'Base',
        component: () => import("@/views/Dashboard.vue"),
        icon: '􀎟',
        children: [
            {
                path: '/',
                name: 'dashboard',
                component: () => import("@/views/Home.vue"),
                icon: '􀎟',
            },
            {
                path: '/apps',
                name: 'apps',
                component: () => import("@/views/apps/Apps.vue"),
                icon: '􀎟',
            }, {
                path: '/beampath',
                name: 'beampath',
                component: () => import("@/views/apps/BeamPath.vue"),
                icon: '􀎟',
            },
            {
                path: '/entity/:entityId',
                name: 'editEntity',
                component: () => import("@/views/home/EntityEdit.vue"),
                icon: '􀎟',
            },
            {
                path: '/extra/homework',
                name: 'homework',
                component: () => import("@/views/extra/Homework.vue"),
                icon: '􀎟',
            },
            {
                path: '/subroutines',
                name: 'subroutines',
                component: () => import("@/views/home/SubRoutineList.vue"),
                icon: '􀎟',
            },
            {
                path: '/sentry',
                name: 'sentry',
                component: () => import("@/views/beam/Beam.vue"),
                icon: '􀎟',
            },
            {
                path: '/monitor',
                name: 'monitor',
                component: () => import("@/views/monitor/MonitorList.vue"),
                icon: '􀎟',
            },
            {
                path: '/lights',
                name: 'lights',
                component: () => import("@/views/apps/Lights.vue"),
                icon: '􀎟',
            },
            {
                path: '/menu',
                name: 'menu',
                component: () => import("@/views/MenuPage.vue"),
                icon: '􀎟',
            },
            {
                path: '/thermostat',
                name: 'thermostat',
                component: () => import("@/views/home/ThermostatView.vue"),
                icon: '􀎟',
            },
            {
                path: '/macros',
                name: 'macros',
                component: () => import("@/views/home/MacroList.vue"),
                icon: '􀎟',
            },
            {
                path: '/apps/request',
                name: 'request',
                component: () => import("@/views/apps/RequestBuilder.vue"),
                icon: '􀎟',
            },
            {
                path: '/apps/message',
                name: 'message',
                component: () => import("@/views/apps/Message.vue"),
                icon: '􀎟',
            },
            {
                path: '/apps/todo',
                name: 'request',
                component: () => import("@/views/apps/Todo.vue"),
                icon: '􀎟',
            },
            {
                path: '/beam',
                name: 'beam',
                component: () => import("@/views/beam/Beam.vue"),
                icon: '􀎟',
            },
            settings,
            {path: "*", component: NotFound}
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
