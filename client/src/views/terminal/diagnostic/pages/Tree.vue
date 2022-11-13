<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Remote} from "@/remote";

const remote = inject("remote") as Remote
let root: TreeNode;

const state = reactive({
    ctx: {} as CanvasRenderingContext2D,
    width: 0,
    height: 0,
    loaded: false,
})

class TreeNode {
    children: TreeNode[]
    name: string

    constructor(name: string) {
        this.name = name
        this.children = []
    }

    addChild(node: TreeNode) {
        this.children.push(node)
    }

    getChildren(): TreeNode[] {
        return this.children
    }

}

onMounted(() => {
    buildTree()
    setupCanvas()
})

watchEffect(() => {
    buildTree()
    redraw()
})

function numKin(node: TreeNode): number {
    return node.children.length
}

function buildTree() {

    let rt = new TreeNode("Root")
    for (let i = 0; i < remote.modules.length; i++) {

        let ents = remote.entities.filter(a => a.module === remote.modules[i].name)

        let modNode = new TreeNode(remote.modules[i].name)
        for (let b = 0; b < ents.length; b++) {

            let tn = new TreeNode(ents[b].name)

            let attrs = remote.attributes.filter(a => a.entity === ents[b].id)
            for (let j = 0; j < attrs.length; j++) {
                if (!attrs[j]) continue
                tn.addChild(new TreeNode(attrs[j].key))
            }
            if (tn.children.length > 0) {
                modNode.addChild(tn)
            }

        }
        if (modNode.children.length > 0) {
            rt.addChild(modNode)

        }

    }

    root = rt
}

function drawTree() {

    drawNode(root, 0, 0, state.width * 1, 60)
}

function drawNode(node: TreeNode, x: number, y: number, w: number, h: number) {
    let ctx = state.ctx

    ctx.fillStyle = "rgba(255,255,255,1)"
    ctx.strokeStyle = "rgba(255,255,255,1)"

    let fontSize = 20
    ctx.font = `500 ${fontSize}px SF Compact`
    let m = ctx.measureText(node.name)
    ctx.fillText(node.name, x + (w / 2) - m.width / 2, y + h / 2 + fontSize / 3)
    ctx.beginPath()
    ctx.rect(x, y, w, h)
    ctx.stroke()
    ctx.closePath()
    let children = node.getChildren()
    let chld = children.sort((a, b) => b.children.length - a.children.length)
    let numChildren = children.length
    for (let i = 0; i < numChildren; i++) {

        drawNode(chld[i], x + (w / numChildren) * i + 2, y + h + 20, (w / numChildren) - 4, h)
    }
}


function redraw() {
    state.ctx.clearRect(0, 0, state.width, state.height)

    state.ctx.rect(0, 0, state.width, state.height)
    state.ctx.stroke()
    drawTree()
}

function setupCanvas() {
    const chart = document.getElementById(`diagnostic-tree`) as HTMLCanvasElement
    if (!chart) return;

    const ctx = chart.getContext('2d')
    if (!ctx) return
    state.ctx = ctx
    let scale = 1.75
    ctx.scale(scale, scale)


    chart.width = chart.clientWidth * scale
    chart.height = chart.clientHeight * scale

    state.width = ctx.canvas.width
    state.height = ctx.canvas.height
    ctx.translate(0, 0)
    ctx.clearRect(0, 0, state.width, state.height)

    redraw()
}

function updateStats() {

    state.loaded = true

}


</script>

<template>
    <div class="page-grid w-100 h-00">
        <canvas id="diagnostic-tree" style="height: 100%; width: 100%;"></canvas>
    </div>
</template>

<style scoped>

.page-grid {
    width: 100%;
    height: 100%;
    display: grid;
    grid-column-gap: 0.25rem;
    grid-row-gap: 0.25rem;
    grid-template-rows: repeat(1, 1fr);
    grid-template-columns: repeat(1, 1fr);
}
</style>