// Copyright (c) 2023 Braden Nicholson

// @ts-ignore
import {registerSW} from 'virtual:pwa-register';

const updateSW = registerSW({
    onNeedRefresh() {
        // show a prompt to the user
    },
    onOfflineReady() {
        // show a ready to work offline to the user
    },
});