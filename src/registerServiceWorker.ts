// Copyright (c) 2023 Braden Nicholson

// @ts-ignore
import {registerSW} from 'virtual:pwa-register'

const updateSW = registerSW({
    onNeedRefresh() {
        console.log("REFRESH")
    },
    onOfflineReady() {
        console.log("OFFLINE")
    },
})

export default updateSW;