// Copyright (c) 2022 Braden Nicholson

// Diagnostic Routes
import Diagnostic from "./Diagnostic.vue";
import Overview from "./pages/Overview.vue";

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
    ]
}