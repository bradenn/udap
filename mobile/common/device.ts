// Copyright (c) 2022 Braden Nicholson

import {Remote} from "./remote";
import {onMounted, reactive} from "vue";
import endpointService from "./services/endpointService";

export interface Device {
    notifications: {
        supported: boolean,
        granted: boolean,
        push: {
            supported: boolean,
            granted: boolean,
            vapidUrl: string,
            requestPermission: () => void
        }
        appBadge: {
            supported: boolean,
            setBadge: (val: number) => void
        }
        requestPermission: () => void
    },

}

function urlBase64ToUint8Array(base64String: string) {
    const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
    const base64 = (base64String + padding)
        .replace(/\-/g, '+')
        .replace(/_/g, '/');
    const rawData = window.atob(base64);
    return Uint8Array.from([...rawData].map(char => char.charCodeAt(0)));
}


export interface Subscription {
    endpoint: string
    expirationTime: any
    keys: Keys
}

export interface Keys {
    p256dh: string
    auth: string
}

function parseJwt(token: string): string {
    let base64Url = token.split('.')[1];
    let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    let jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
}

export default function useDevice(remote: Remote): Device {

    const state = reactive({
        notifications: {
            supported: false,
            granted: false,
            push: {
                supported: false,
                granted: false,
                vapidUrl: "Unknown",
                requestPermission: requestPush
            },
            appBadge: {
                supported: false,
                setBadge: setBadge
            },
            requestPermission: requestPermission
        },
    } as Device)

    onMounted(() => {
        updateNotificationStatus()
    })

    function requestPush() {
        navigator.serviceWorker.ready.then((registration) => {
            registration.showNotification("Notification", {
                body: "Hello there bozo!",
                icon: "",
                tag: "vibration-sample",
            });
        });
    }


    function pushSubscribe() {
        navigator.serviceWorker.ready.then(function (registration: ServiceWorkerRegistration) {

            const vapidPublicKey = 'BN_adwCB27Su3AexUbrWkV1Gzdhfgtuh7CRGLN_aw_O0PvHdv6SrLXGO9bIxSs8DjZDy7BjwvFAsWuGarILDmSU';

            return registration.pushManager.subscribe({
                userVisibleOnly: true,
                applicationServerKey: urlBase64ToUint8Array(vapidPublicKey),
            });

        })
            .then(function (subscription) {

                state.notifications.push.granted = !(!subscription);

            })
            .catch(err => console.error(err));
    }


    function verify() {

        state.notifications.push.vapidUrl = "Unset"
        if ('serviceWorker' in navigator) {
            navigator.serviceWorker.ready.then(function (registration) {
                return registration.pushManager.getSubscription();
            }).then(function (subscription) {
                if (!subscription) {
                    state.notifications.push.vapidUrl = "Unset"
                    pushSubscribe();
                } else {
                    state.notifications.push.vapidUrl = subscription.endpoint
                    state.notifications.push.granted = true
                    let token = localStorage.getItem("token");
                    if (token) {
                        let jwt = JSON.parse(parseJwt(token)) as { id: string }
                        endpointService.registerPush(jwt.id, JSON.stringify(subscription)).then(r => console.log(`Registered Endpoint: ${jwt.id}`))
                    }
                }
            });
        }
    }

    function setBadge(value: number): void {
        if ("setAppBadge" in navigator) {
            state.notifications.appBadge.supported = true
            // const n = new Notification(`Hi! Number is set to ${value}!!`, {
            //     tag: "soManyCuteNotification",
            // });
            //@ts-ignore
            navigator.setAppBadge(Math.round(value))
        }
    }

    function requestPermission(): void {
        if ("Notification" in window) {
            state.notifications.supported = true
            Notification.requestPermission().then((permission) => {
                state.notifications.granted = permission === 'granted';
            }).catch(err => alert(err))
        } else {
            state.notifications.supported = false
        }
        updateNotificationStatus()
    }

    function updateNotificationStatus() {
        if ("setAppBadge" in navigator) {
            state.notifications.appBadge.supported = true
        }

        if ("Notification" in window) {
            state.notifications.supported = true
            state.notifications.push.supported = true
            state.notifications.granted = Notification.permission === 'granted';
            verify()
        } else {
            state.notifications.supported = false
        }

    }

    return state
}

