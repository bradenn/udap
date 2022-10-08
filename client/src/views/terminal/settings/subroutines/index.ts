// Copyright (c) 2022 Braden Nicholson

import SubroutineView from "@/views/terminal/settings/subroutines/SubroutineView.vue";

// Overview
import Subroutines from "@/views/terminal/settings/subroutines/Subroutines.vue";
import Create from "@/views/terminal/settings/subroutines/pages/Create.vue";

const routes = {
    path: '/terminal/settings/subroutines',
    name: 'Subroutines',
    icon: 'timeline',
    meta: {
        order: 3,
    },
    children: [
        {
            path: '/terminal/settings/subroutines',
            name: 'Subroutines',
            icon: 'timeline',
            component: Subroutines
        },
        {
            path: '/terminal/settings/subroutines/create',
            name: 'Create A Subroutine',
            icon: 'timeline',
            component: Create
        }
    ],
    component: SubroutineView
}

export default routes