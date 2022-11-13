// Copyright (c) 2022 Braden Nicholson

import SubroutineView from "@/views/terminal/settings/subroutines/SubroutineView.vue";

// Overview
import Subroutines from "@/views/terminal/settings/subroutines/Subroutines.vue";
import Create from "@/views/terminal/settings/subroutines/pages/Create.vue";
import EditSubroutine from "@/views/terminal/settings/subroutines/pages/EditSubroutine.vue";
import Trigger from "@/views/terminal/settings/subroutines/pages/NewTrigger.vue";
import Macro from "@/views/terminal/settings/subroutines/pages/CreateMacro.vue";
import Zone from "@/views/terminal/settings/subroutines/pages/CreateZone.vue";

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
        }, {
            path: '/terminal/settings/subroutines/macros/create',
            name: 'Create A Macro',
            icon: 'timeline',
            component: Macro
        },
        {
            path: '/terminal/settings/subroutines/macros/:macroId/edit',
            name: 'Edit A Macro',
            icon: 'timeline',
            component: Macro
        },
        {
            path: '/terminal/settings/subroutines/trigger',
            name: 'Create A Trigger',
            icon: 'switch',
            component: Trigger
        },
        {
            path: '/terminal/settings/subroutines/zones/create',
            name: 'Create A Zone',
            icon: 'switch',
            component: Zone
        },
        {
            path: '/terminal/settings/subroutines/:subroutine/edit',
            name: 'Edit Subroutine',
            icon: 'switch',
            component: EditSubroutine
        }
    ],
    component: SubroutineView
}

export default routes