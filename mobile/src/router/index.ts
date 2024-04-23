// Copyright (c) 2022 Braden Nicholson

import {createRouter, createWebHistory} from "vue-router";

import setupRoutes from "@/views/setup";


import settings from "@/views/settings";
import apps from "@/views/apps";
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
                path: '/menu',
                name: 'menu',
                component: () => import("@/views/MenuPage.vue"),
                icon: '􀎟',
            },

            {
                path: '/macros',
                name: 'macros',
                component: () => import("@/views/home/MacroList.vue"),
                icon: '􀎟',
            },

            apps,
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
