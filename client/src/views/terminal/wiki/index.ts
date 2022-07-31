// Copyright (c) 2022 Braden Nicholson

// Wiki Routes
import Wiki from "./Wiki.vue";
import WikiMenu from "./WikiMenu.vue";

const wikiRoutes = {
    path: '/terminal/wiki',
    name: 'Wiki',
    component: Wiki,
    icon: 'fa-book',
    children: [
        {
            path: '/terminal/wiki',
            name: 'WikiMenu',
            component: WikiMenu,
        },
    ]
}

export default wikiRoutes