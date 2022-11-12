// Copyright (c) 2022 Braden Nicholson


import {reactive} from "vue";
import {v4 as uuidv4} from "uuid";

interface Notifications {
    title: string
    message: string
    severity: number
    payload: string
    timeout: number
}

interface ToastObject {
    name: string,
    message: string,
    severity: number,
    priority: number,
    duration: number,
    uuid: string,
}

let toasts = reactive<{
    queue: ToastObject[];
    active: boolean;
    interval: number;
    current: ToastObject;
}>({
    queue: [] as ToastObject[],
    active: false,
    interval: 0,
    current: {} as ToastObject
})


function updateQueue() {
    if (toasts.queue.length > 0 && !toasts.active) {
        toasts.active = true
        toasts.current = toasts.queue[0]

        if (toasts.interval != 0) return
        toasts.interval = setInterval(() => {
            if (toasts.active) {
                toasts.current.duration -= 16
                if (toasts.current.duration <= 0) {
                    toasts.queue = toasts.queue.filter(t => t.uuid != toasts.current.uuid)
                    if (toasts.queue.length < 1) {
                        toasts.active = false
                        clearInterval(toasts.interval)
                        toasts.interval = 0
                    } else {
                        toasts.current = toasts.queue[0]
                    }
                }
            }
        }, 16)
    }
}


function success() {

    toasts.queue.push({
        name: "System",
        message: "Operation successful",
        severity: 1,
        priority: 1,
        duration: 2400,
        uuid: uuidv4()
    })

    updateQueue()
}

function show(name: string, message: string, severity: number, duration: number) {
    toasts.queue.push({
        priority: 0,
        name: name,
        message: message,
        severity: severity,
        duration: duration,
        uuid: uuidv4()
    })
    updateQueue()
}

function getQueue(): ToastObject[] {
    return toasts.queue
}


export default {
    success,
    show,
    toasts
}


