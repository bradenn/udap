import General from "@/views/settings/General.vue";
import Connection from "@/views/settings/Connection.vue";
import SettingsMenu from "@/views/settings/SettingsMenu.vue";
import Background from "@/views/settings/Background.vue";
import Color from "@/views/settings/Color.vue";
import Entities from "@/views/settings/entities/Entities.vue";
import EntityManage from "@/views/settings/entities/EntityManage.vue";
import Frame from "@/views/settings/entities/Frame.vue";
import Config from "@/views/settings/Config.vue";
import Modules from "@/views/settings/Modules.vue";
import Zones from "@/views/settings/Zones.vue";
import Endpoints from "@/views/settings/Endpoints.vue";
import EditHome from "@/views/settings/EditHome.vue";
import Timings from "@/views/settings/Timings.vue";
import WebServices from "@/views/settings/WebServices.vue";
import Module from "@/views/settings/module/Module.vue";


export default {
    path: '/home/settings',
    name: 'settings',
    component: () => import('./Settings.vue'),
    children: [
        {
            path: '/home/settings',
            name: 'menu',
            component: SettingsMenu,
        },
        {
            path: '/home/settings/general',
            name: 'general',
            component: General,
        },
        {
            path: '/home/settings/connection',
            name: 'connection',
            component: Connection,
        },
        {
            path: '/home/settings/background',
            name: 'background',
            component: Background,
        },
        {
            path: '/home/settings/config',
            name: 'configuration',
            component: Config,
        },
        {
            path: '/home/settings/modules',
            name: 'modules',
            component: Modules,
            children: []
        },
        {
            path: '/home/settings/modules/:moduleId',
            name: 'module',
            component: Module,
        },
        {
            path: '/home/settings/color',
            name: 'color',
            component: Color,
        },
        {
            path: '/home/settings/zones',
            name: 'zones',
            component: Zones,
        },
        {
            path: '/home/settings/home',
            name: 'home',
            component: EditHome,
        },
        {
            path: '/home/settings/timings',
            name: 'timings',
            component: Timings,
        },
        {
            path: '/home/settings/webservices',
            name: 'webservices',
            component: WebServices,
        },
        {
            path: '/home/settings/endpoints',
            name: 'endpoints',
            component: Endpoints,
        },
        {
            path: '/home/settings/entities',
            name: 'entity_frame',
            component: Frame,
            children: [
                {
                    path: '/home/settings/entities',
                    name: 'Entities',
                    component: Entities,
                },
                {
                    path: '/home/settings/entities/:entityId',
                    name: 'entities_view',
                    component: EntityManage,
                }
            ]
        }
    ]
}
