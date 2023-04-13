// Copyright (c) 2022 Braden Nicholson

import {createRouter, createWebHashHistory} from "vue-router";
import Dashboard from "@/views/Dashboard.vue";
import setupRoutes from "@/views/setup";
import SubRoutineList from "@/views/home/SubRoutineList.vue";
import MacroList from "@/views/home/MacroList.vue";
import Home from "@/views/Home.vue";
import EntityEdit from "@/views/home/EntityEdit.vue";

const mobileRoutes =
    {
        path: '/home',
        name: 'Home',
        redirect: "/home/dashboard",
        component: Dashboard,
        icon: '􀎟',
        children: [
            {
                path: '/home/dashboard',
                name: 'dashboard',
                component: Home,
                icon: '􀎟',
            },
            {
                path: '/home/entity/:entityId',
                name: 'editEntity',
                component: EntityEdit,
                icon: '􀎟',
            },
            {
                path: '/home/subroutines',
                name: 'subroutines',
                component: SubRoutineList,
                icon: '􀎟',
            },
            {
                path: '/home/macros',
                name: 'macros',
                component: MacroList,
                icon: '􀎟',
            },
        ]
    };

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
