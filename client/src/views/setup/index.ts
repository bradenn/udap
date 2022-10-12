// Copyright (c) 2022 Braden Nicholson

import Setup from "@/views/setup/Setup.vue";
import Authentication from "@/views/setup/Authentication.vue";
import Welcome from "@/views/setup/Controller.vue";

export default {
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
};