// Copyright (c) 2023 Braden Nicholson
//
// //@ts-ignore
// import {registerSW} from "virtual:pwa-register";


self.addEventListener('push', (event: PushEvent) => {
    console.log("Hello bozo! (From SW)")


    const request = JSON.parse(event?.data.text()) as {
        endpointId: string,
        message: string,
        title: string
    };
    const icon = "/pwa196x196.png";
    event.waitUntil(
        self.registration.showNotification(request.title, {
            body: `${request.message}`,
            tag: `${request.endpointId}.${request.title}`,
            icon,
        })
    )
});


// registerSW({
//     onRegisteredSW: registered
// })()


