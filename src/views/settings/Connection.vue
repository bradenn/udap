<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Table from "@/components/Table.vue";
import TableText from "@/components/TableText.vue";
import core from "@/core";
import {computed, onBeforeMount, onMounted, reactive} from "vue";
import type {Endpoint, Timing} from "@/types";


const remote = core.remote()

const state = reactive({
  endpoint: {} as Endpoint,
  storage: {},
  ctx: {} as CanvasRenderingContext2D
})

let token = localStorage.getItem("token")

interface TokenData {
    id: string
}

function parseJwt(token: string): TokenData {
    let base64Url = token.split('.')[1];
    let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    let jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

  return JSON.parse(jsonPayload) as TokenData;
}

onBeforeMount(() => {
  let tokenMeta = parseJwt(token)
  let ep = remote.endpoints.find(e => e.id === tokenMeta.id)
  if (!ep) return
  state.endpoint = ep
  state.storage = localStorage.getItem("preferences")

})

onMounted(() => {
  setup()
  animate()
})

function animate() {
  draw()
  requestAnimationFrame(animate)
}

let timings = computed(() => {
  return correlateTimings(remote.timings)
})

function draw() {
  let ctx = state.ctx
  let values = timings.value
  let w = ctx.canvas.width
  let h = ctx.canvas.height

  ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height)
  ctx.strokeStyle = "rgba(255,128,12,0.8)"
  ctx.beginPath()
  ctx.moveTo(0, h / 2)
  let dx = w / (values.length - 1)
  for (let i = 0; i < values.length; i++) {
    ctx.lineTo(i * dx, (h / 2) - values[i] * h / 4)
  }
  ctx.moveTo(w, h / 2)

  ctx.closePath()

  ctx.stroke()

}


function correlateTimings(t: Timing[]): number[] {
  let out = []
  for (let i = 0; i < t.length; i++) {
    let deltaNano = t[i].stopNano - t[i].startNano
    let b = Math.log(deltaNano)
    out.push(Math.round(b))
  }
  let mx = Math.max(...out)
  let nx = Math.min(...out)
  out = out.map(o => (o / mx))

  return out
}


function setup() {
  let canvas = document.getElementById("performance") as HTMLCanvasElement
  if (!canvas) return
  let ctx = canvas.getContext("2d") as CanvasRenderingContext2D
  if (!ctx) return
  state.ctx = ctx
  ctx.canvas.width = canvas.clientWidth * 2
  ctx.canvas.height = canvas.clientHeight * 2

}


</script>

<template>
    <div class="d-flex flex-column gap-3">
        <Table title="Client">
            <TableText icon="" title="Status">{{ remote.connected ? "Connected" : "Disconnected" }}</TableText>
            <TableText icon="" title="Endpoint Identifier">{{ state.endpoint.name }}</TableText>
            <TableText icon="" title="Permissions">
              <div v-if="state.endpoint.type === 'terminal'">
                Inexorable
              </div>
              <div v-else>

              </div>
            </TableText>
          <TableText icon="" title="Local Storage">{{ state.storage }}</TableText>
        </Table>

      <div class="element">
        <canvas id="performance" style="width: 100%; height: 4rem;"></canvas>
      </div>

      <Table v-if="remote.metadata.system" title="Node">
        <TableText icon="" title="Hostname">{{ remote.metadata.system.hostname }}</TableText>
        <TableText icon="" title="IPv4">{{ remote.metadata.system.ipv4 }}</TableText>
        <TableText icon="" title="Cores">{{ remote.metadata.system.cores }}</TableText>
        <TableText icon="" title="Environment">{{ remote.metadata.system.environment }}</TableText>
        <TableText icon="" title="Version">{{ remote.metadata.system.version }}</TableText>
        <TableText icon="" title="Compiler">{{ remote.metadata.system.go }}</TableText>
      </Table>
      <div v-else> Disconnected from Network</div>
    </div>
</template>

<style scoped>

</style>