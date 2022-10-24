// Copyright (c) 2022 Braden Nicholson

// Diagnostic Routes
import Diagnostic from "./Diagnostic.vue";
import Overview from "./pages/Overview.vue";
import Tree from "./pages/Tree.vue";
import Routes from "./pages/Routes.vue";

export default {
    path: '/terminal/diagnostics',
    name: 'Diagnostics',
    redirect: '/terminal/diagnostics/summary',
    component: Diagnostic,
    icon: 'fa-tools',
    meta: {
        title: "Diagnostics",
        status: "wip"
    },
    children: [
        {
            path: '/terminal/diagnostics/summary',
            name: 'DiagnosticOverview',
            meta: {
                title: "Overview"
            },
            component: Overview,
        },
        {
            path: '/terminal/diagnostics/tree',
            name: 'DiagnosticTree',
            icon: 'tools',
            meta: {
                title: "Tree"
            },
            component: Tree,
        },
        {
            path: '/terminal/diagnostics/routes',
            name: 'Routes',
            icon: 'tools',
            meta: {
                title: "Routes"
            },
            component: Routes,
        },
    ]
}