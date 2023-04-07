// Copyright (c) 2022 Braden Nicholson

import Enroll from "@/views/setup/Enroll.vue";
import Setup from "@/views/setup/Setup.vue";

const setupRoutes = {
    path: '/setup',
    name: 'setup',
    redirect: '/setup/enroll',
    component: Setup,
    children: [
        {
            path: '/setup/enroll',
            name: 'Enroll',
            component: Enroll,
            icon: '􀎟',
        },
    ],
}

export default setupRoutes
