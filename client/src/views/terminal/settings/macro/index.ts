// Copyright (c) 2022 Braden Nicholson

import MacroList from "@/views/terminal/settings/macro/MacroList.vue";
import Macros from "@/views/terminal/settings/macro/Macros.vue";


const routes = {
    path: '/terminal/settings/macros',
    name: 'Macros',
    component: Macros,
    meta: {
        order: 3,
        icon: '􀩶'
    },
    icon: '􀍟',
    children: [
        {
            path: '/terminal/settings/macros',
            name: 'MacrosList',
            icon: '􀐗',
            component: MacroList
        }
    ]
}

export default routes