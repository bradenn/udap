// Copyright (c) 2022 Braden Nicholson

import ModuleView from "@/views/terminal/settings/module/ModuleView.vue";

// Overview
import Modules from "@/views/terminal/settings/module/Modules.vue";

// Editor
import Module from "@/views/terminal/settings/module/Module.vue";

// Pages
import Metadata from "@/views/terminal/settings/module/pages/Metadata.vue";
import Versions from "@/views/terminal/settings/module/pages/Versions.vue";
import Actions from "@/views/terminal/settings/module/pages/Actions.vue";
import Config from "@/views/terminal/settings/module/pages/Config.vue";

const routes = {
    path: '/terminal/settings/modules',
    name: 'Modules',
    icon: 'layer-group',
    meta: {
        order: 2,
    },
    children: [
        {
            path: '/terminal/settings/modules',
            name: 'Module Overview',
            icon: 'layer-group',
            component: Modules
        },
        {
            path: '/terminal/settings/modules/:moduleId',
            redirect: (to: any) => ({
                name: "Module Actions",
                query: {moduleId: to.params.moduleId},
            }),
            name: 'Module',
            component: Module,
            children: [
                {
                    path: '/terminal/settings/modules/:moduleId/actions',
                    name: 'Module Actions',
                    component: Actions,
                },
                {
                    path: '/terminal/settings/modules/:moduleId/config',
                    name: 'Module Config',
                    component: Config,
                },
                {
                    path: '/terminal/settings/modules/:moduleId/metadata',
                    name: 'Module Metadata',
                    component: Metadata,
                },
                {
                    path: '/terminal/settings/modules/:moduleId/versions',
                    name: 'Module Versions',
                    component: Versions,
                }
            ]
        }
    ],
    component: ModuleView
}

export default routes