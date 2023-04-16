import Settings from "@/views/settings/Settings.vue";
import General from "@/views/settings/General.vue";
import Connection from "@/views/settings/Connection.vue";
import SettingsMenu from "@/views/settings/SettingsMenu.vue";
import Background from "@/views/settings/Background.vue";
import Color from "@/views/settings/Color.vue";


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
            path: '/home/settings/color',
            name: 'color',
            component: Color,
        }
    ]
}
