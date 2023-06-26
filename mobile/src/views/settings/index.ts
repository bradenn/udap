import Settings from "@/views/settings/Settings.vue";
import General from "@/views/settings/General.vue";
import Connection from "@/views/settings/Connection.vue";
import SettingsMenu from "@/views/settings/SettingsMenu.vue";
import Background from "@/views/settings/Background.vue";
import Color from "@/views/settings/Color.vue";
import Entities from "@/views/settings/entities/Entities.vue";
import EntityManage from "@/views/settings/entities/EntityManage.vue";
import Frame from "@/views/settings/entities/Frame.vue";
import Config from "@/views/settings/Config.vue";


export default {
    path: '/home/settings',
    name: 'settings',
    component: Settings,
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
            path: '/home/settings/color',
            name: 'color',
            component: Color,
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
