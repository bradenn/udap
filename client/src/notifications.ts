// Copyright (c) 2022 Braden Nicholson

import {reactive} from "vue";
import {v4 as uuidv4} from "uuid";

export interface NotifyBody {
    name: string,
    message: string,
    trace?: string,
    severity: number,
    priority: number,
    duration: number,
    uuid: string,
}

interface NotifyState {
    log: NotifyBody[];
    queue: NotifyBody[];
    active: boolean;
    interval: number;
    current: NotifyBody;
    currentWait: number;
    begin: number;
}

export interface Notify {
    success(name: string, message: string): void

    fail(name: string, message: string, trace?: string): void

    show(name: string, message: string, severity: number, duration: number): void

    isActive(): boolean

    current(): NotifyBody

    logs(): NotifyBody[]

    clearLog(): void

}

let state = reactive<NotifyState>({
    log: [] as NotifyBody[],
    queue: [] as NotifyBody[],
    active: false,
    interval: 0,
    current: {} as NotifyBody,
    begin: 0,
    currentWait: 0,
})


function logs(): NotifyBody[] {
    return state.log
}

// Clears the logged notifications
function clearLog() {
    state.log = state.log.filter(_ => false)
}

// Promotes the first element in the queue to the new current
function dequeue() {

    // Pop an element from the front of the queue
    let next = state.queue.shift()
    // Make sure the element is not null
    if (!next) return
    // Update the current element
    state.current = next
    state.log.push(next)
    // Update the start time
    state.begin = Date.now().valueOf()
    // Update the countdown duration
    state.currentWait = next.duration
}

// Begins the countdown for the current element
function countdown() {
    // Request the next animation frame
    state.interval = requestAnimationFrame(countdown)
    // Check if the queue is currently active
    if (state.active) {
        // Set the current elapsed duration to the time since the beginning minus the original timeout
        state.current.duration = state.currentWait - (Date.now().valueOf() - state.begin)
        // Check if the duration has fully elapsed
        if (state.current.duration <= 0) {
            // Check if the queue is empty
            if (state.queue.length <= 0) {
                // Disable the notification queue
                state.active = false
                // Cancel the animation from to stop updated
                cancelAnimationFrame(state.interval)
                // Reset the animation frame variable
                state.interval = 0
            } else {
                // Promote the next node to be the current
                dequeue()
            }
        }
    }
}

// Forces an update on the notifications queue
function updateQueue() {
    // Check if the queue is not empty and the queue is not active (not is use)
    if (state.queue.length > 0 && !state.active) {
        // Enable the queue
        state.active = true
        // Dequeue the next in row to be the current
        dequeue()
        // Return now if the countdown is already going
        if (state.interval != 0) return
        // Begin countdown
        countdown()
    }
}

// Push a notification into the viewing queue
function pushQueue(toast: NotifyBody) {

    state.queue.push(toast)

    updateQueue()
}

// Sends a success notification to the user
function success(name: string, message: string): void {
    pushQueue({
        name: name,
        message: message,
        severity: 1,
        priority: 1,
        duration: 3000,
        uuid: uuidv4()
    })
}

// Sends a failure notification to the user
function fail(name: string, message: string, trace?: string): void {
    pushQueue({
        name: name,
        message: message,
        severity: 3,
        priority: 2,
        duration: 3000,
        trace: trace,
        uuid: uuidv4()
    })
}

// Send a custom message to the user
function show(name: string, message: string, severity: number, duration: number) {
    pushQueue({
        priority: 0,
        name: name,
        message: message,
        severity: severity,
        duration: duration,
        uuid: uuidv4()
    })
}


// Send a custom message to the user
function current(): NotifyBody {
    return state.current
}

// Send a custom message to the user
function isActive(): boolean {
    return state.active
}

export default {
    success,
    fail,
    show,
    current,
    isActive,
    logs,
    clearLog
}



