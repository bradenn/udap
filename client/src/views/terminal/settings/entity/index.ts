// Copyright (c) 2022 Braden Nicholson

import Index from "@/views/terminal/settings/entity/pages/Index.vue";
import Entities from "@/views/terminal/settings/entity/Entities.vue";
import Entity from "@/views/terminal/settings/entity/pages/Entity.vue";

const routes = {
    path: '/terminal/settings/entities',
    name: 'Entities',
    icon: '􀩶',
    meta: {
        order: 3,
        icon: '􀩶'
    },
    children: [
        {
            path: '/terminal/settings/entities',
            name: 'EntityIndex',
            icon: '􀩶',
            component: Index
        },
        {
            path: '/terminal/settings/entities/:entityId',
            name: 'EditEntity',
            icon: '􀩶',
            component: Entity
        },
    ],
    component: Entities
}

export default routes