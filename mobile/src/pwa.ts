// Copyright (c) 2023 Braden Nicholson
//
//@ts-ignore
// import {registerSW} from "virtual:pwa-register";


declare let self: ServiceWorkerGlobalScope
//
import {cleanupOutdatedCaches, precacheAndRoute} from 'workbox-precaching'

cleanupOutdatedCaches()

precacheAndRoute(self.__WB_MANIFEST || [])
// Check for support first.
//@ts-ignore
if (navigator.setAppBadge) {
    // Display the number of unread messages.
    //@ts-ignore
    navigator.setAppBadge(0).then(r => console.log(r)).catch(r => console.log(r));
}

self.addEventListener('push', (event: PushEvent) => {

    if (!event) {
        return
    }

    if (!event.data) {
        return;
    }

    const request = JSON.parse(event.data.text()) as {
        endpointId: string,
        message: string,
        title: string
    };


    const icon = "/pwa196x196.png";
    //@ts-ignore
    event.waitUntil(
        //@ts-ignore
        self.registration.showNotification(request.title, {
            body: `${request.message}`,
            tag: `${request.endpointId}.${request.title}`,
            icon,
        })
    )
});
export {}
//
// registerSW({
//     onRegisteredSW: () => {
//         console.log("anal")
//     }
// })()


