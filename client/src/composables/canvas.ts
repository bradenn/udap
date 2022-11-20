// Copyright (c) 2022 Braden Nicholson

import {reactive} from "vue";
import {v4 as uuidv4} from "uuid";

let state = reactive({
    animationFrame: 0,
    width: 0,
    height: 0,
    ctx: {} as CanvasRenderingContext2D,
    uuid: uuidv4(),
    drawFunc: defaultDraw as (c: CanvasRenderingContext2D) => void,
})

function animate() {
    state.animationFrame = requestAnimationFrame(animate)
    redraw()
}

function defaultDraw(c: CanvasRenderingContext2D) {
    let levels = 40

    let dx = Math.ceil(state.height / levels)

    for (let i = 0; i < levels; i++) {
        c.strokeRect(0, 0, state.width, dx * i)

    }

}

function redraw() {
    state.ctx.clearRect(0, 0, state.width, state.height)
    state.drawFunc(state.ctx)
}


function dispose() {
    cancelAnimationFrame(state.animationFrame)
    state.ctx.canvas.remove()

}

function setupCanvas(draw: (c: CanvasRenderingContext2D) => void) {
    const chart = document.getElementById(`${state.uuid}`) as HTMLCanvasElement
    if (!chart) return;

    const ctx = chart.getContext('2d')
    if (!ctx) return

    state.ctx = ctx
    let scale = 2
    ctx.scale(scale, scale)

    if (draw) {
        state.drawFunc = draw
    }
    chart.width = chart.clientWidth * scale
    chart.height = chart.clientHeight * scale

    ctx.translate(0.5, 0.5)

    state.width = ctx.canvas.width
    state.height = ctx.canvas.height

    ctx.clearRect(0, 0, state.width, state.height)


    state.animationFrame = requestAnimationFrame(animate)
}

export default {
    setupCanvas,
    dispose,
    uuid: state.uuid,
}

