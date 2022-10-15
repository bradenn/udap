// Copyright (c) 2022 Braden Nicholson


interface HapticPulse {
    freq: number,
    power: number,
    amplitude: number,
}

class HapticEngine {
    url: string
    ws: WebSocket
    ready: boolean

    constructor(url: string) {
        this.url = url
        this.ws = this.connect()

        this.ready = false
    }

    public tick() {
        this.sendPulse({
            freq: 1,
            power: 1,
            amplitude: 2048
        })
    }

    public tap(frequency: number, iterations: number, amplitude: number) {
        this.sendPulse({
            freq: frequency,
            power: iterations,
            amplitude: this.map_range(amplitude, 0, 100, 0, 4095)
        })
    }

    private connect(): WebSocket {
        let ws = new WebSocket(this.url)
        ws.addEventListener("open", (_: Event) => {
            this.ready = true
            console.log("CONNECTED")
        })
        ws.addEventListener("close", (_: CloseEvent) => {
            this.ready = false
            // Try to reconnect after 5 seconds
            setTimeout(this.connect, 5000)
        })
        return ws
    }

    private sendPulse(pulse: HapticPulse) {
        if (this.ready) {
            let payload = JSON.stringify(pulse)
            this.ws.send(payload)
        }
    }

    private map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
        return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
    }
}

export {
    HapticEngine
}

