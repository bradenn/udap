// Copyright (c) 2022 Braden Nicholson

import {createRouter, createWebHashHistory} from "vue-router";
import Mobile from "@/views/Mobile.vue";
import Dashboard from "@/views/Dashboard.vue";
import setupRoutes from "@/views/setup";

const mobileRoutes = {
    path: '/mobile',
    name: 'mobile',
    redirect: '/mobile/home',
    component: Mobile,
    children: [
        {
            path: '/mobile/home',
            name: 'Home',
            component: Dashboard,
            icon: '􀎟',
        },
    ],
}

const routes = [
    mobileRoutes,
    setupRoutes
    // {path: '/:pathMatch(.*)*', redirect: "/"}
]

const router = createRouter({
    history: createWebHashHistory(),
    routes: routes,
})

export default router
