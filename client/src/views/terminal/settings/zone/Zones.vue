<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, onUnmounted, reactive, watchEffect} from "vue";
import type {Entity, Zone} from "@/types";
import Loader from "@/components/Loader.vue";
import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import CreateZone from "@/views/terminal/settings/zone/CreateZone.vue";
import axios from "axios";
import FixedScroll from "@/components/scroll/FixedScroll.vue";
import type {Remote} from "@/remote";
import Toolbar from "@/components/toolbar/Toolbar.vue";
import ToolbarButton from "@/components/ToolbarButton.vue";
import core from "@/core";
import zoneService from "@/services/zoneService";

let remote = inject("remote") as Remote
let preferences = inject('preferences')

let state = reactive({
    zones: {} as Zone[],
    entities: [] as Entity[],
    deleted: {} as Zone[],
    selected: "",
    selectedZone: {} as Zone,
  showDeleted: false,
  animationFrame: 0,
  loading: true,
  mode: "list",
  listMode: "view",
  model: "",
  ctx: {} as CanvasRenderingContext2D,
  width: 0,
  height: 0,
  runner: 0,
  tick: 0
})

onMounted(() => {
  state.loading = true
  handleUpdates()
  state.selectedZone = state.zones[0]
  setupCanvas()
  state.loading = false
})


onUnmounted(() => {
  state.ctx.canvas.remove()
  cancelAnimationFrame(state.animationFrame)
})


watchEffect(() => {
  handleUpdates()
  return state.zones
})

function sortZones(a: Zone, b: Zone): number {
  if (a.pinned && !b.pinned) {
    return -1;
  } else if (!a.pinned && b.pinned) {
    return 1;
  } else if (a.pinned && b.pinned && a.name > b.name) {
    return 2;
  } else if (a.pinned && b.pinned && a.name < b.name) {
    return -2;
  }
  return 0
}

function handleUpdates() {
  state.zones = remote.zones.sort((a, b) => sortZones(a, b))
  state.entities = remote.entities

}

function setMode(mode: string) {
  state.mode = mode
}

function selectZone(id: string) {
    if (state.selected === id) {
        state.selected = ""
        state.selectedZone = {} as Zone
        return;
    }
    let target = state.zones.find(z => z.id === id)
    if (!target) return
    state.selectedZone = target
    state.selected = state.selectedZone.id
}

function toggleShowDeleted() {
    state.showDeleted = !state.showDeleted
}

function setListMode(value: string) {
    state.listMode = value
}

const notify = core.notify()

function deleteZone() {
    let local = state.selectedZone
    if (!local) return
    zoneService.setDeleted(local.id, true)
        .then((r) => {
            state.selectedZone.deleted = true
            notify.success("Zone Deleted", `Zone '${local.name}' has been deleted.`)

        })
      .catch((r) => {
          notify.fail("Zone Deleted", `Zone '${local.name}' has not been deleted. `, r)
          console.log(r)
      })
}

function restoreZone() {
  let local = state.selectedZone
  if (!local) return
  axios.post(`http://10.0.1.2:3020/zones/${local.id}/restore`).then((r) => {
    state.selectedZone.deleted = false
  }).catch((r) => {
    console.log(r)
  })
}

function pinZone() {
  let local = state.selectedZone
  if (!local) return
  axios.post(`http://10.0.1.2:3020/zones/${local.id}/pin`)
      .then((r) => {
        state.selectedZone.pinned = true
      })
      .catch((r) => {
        console.log(r)
      })
}

function unpinZone() {
  let local = state.selectedZone
  if (!local) return
  axios.post(`http://10.0.1.2:3020/zones/${local.id}/unpin`).then((r) => {
    state.selectedZone.pinned = false
  }).catch((r) => {
    console.log(r)
  })
}

function toggleEntity(entity: Entity) {
  let target = state.selectedZone.entities.find(e => e.id === entity.id)
  if (!target) {
    state.selectedZone.entities.push(entity)
    axios.post(`http://locahost:3020/zones/${state.selectedZone.id}/entities/${entity.id}/add`).then((r) => {
    }).catch((r) => {
    })
  } else {
    state.selectedZone.entities = state.selectedZone.entities.filter(e => e.id !== entity.id)
    axios.post(`http://locahost:3020/zones/${state.selectedZone.id}/entities/${entity.id}/remove`).then((r) => {
    }).catch((r) => {
    })
  }


}

const livingRoom: number[][] = [
  [0.0001, 0], // Top Left
  [3.66, 0], // Top Right
  [3.66, 2.56],// Right Corner
  [3.2, 3.02], // Right Short wall
  [5.363, 5.183], // Right Long wall
  [4.224, 6.32127], // Front Door
  [3.707, 5.805], // Pantry Door Side
  [3.021, 6.491], // Pantry Face Side
  [3.537, 7.007], // Pantry Fridge Side
  [1.826, 8.7181], // Kitchen Far wall
  [-0.3374, 6.554], // Kitchen Left wall
  [0.285, 5.9317], // Kitchen Laundry Wall
  [-1.115, 4.5316], // Laundry Wall
  [0.3133, 3.103], // Room Wall
  [0.0001, 2.794], // Left Short Wall
]

const bedroom: number[][] = [[0.00007071067811865477, 0.00007071067811865474], [2.0364675298172563, -2.036467529817257], [2.0364675298172563, -2.3764444702117493], [3.1364428386310497, -2.3764444702117498], [3.1364428386310506, 1.993616858877351], [4.440892098500626e-16, 2.0199212311374914]]

function drawZone() {

  state.ctx.lineWidth = 3
  let cx = (state.width / 2) - 160
  let cy = (state.height / 2)
  let scale = 200;

  let ctx = state.ctx
  state.ctx.strokeStyle = "rgba(255,255,255,0.1)"
  ctx.strokeRect(0, 0, state.width, state.height)
  ctx.beginPath()
  ctx.moveTo(state.width / 2, 0)
  ctx.lineTo(state.width / 2, state.height)
  ctx.closePath()
  ctx.stroke()

  ctx.beginPath()
  ctx.moveTo(0, state.height / 2)
  ctx.lineTo(state.width, state.height / 2)
    ctx.closePath()
    ctx.stroke()

    state.ctx.strokeStyle = "rgba(255,149,0,0.6)"
    state.ctx.fillStyle = "rgba(255,255,255,0.025)"

    drawRoom(bedroom, scale)

    state.ctx.fillStyle = "rgba(255,255,255,0.25)"
    ctx.font = "500 32px SF Pro Display"
    if (!state.selectedZone) return
    ctx.fillText(state.selectedZone?.name || "", 10, 36)

    for (let i = 0; i < state.selectedZone.entities.length; i++) {
        let ent = state.selectedZone.entities[i]
        ctx.fillText(ent.name, 10, 36 + 50 + 30 * i)
        if (ent.position != "{}") {
            ctx.fillText(ent.position, 200, 36 + 50 + 30 * i)
            let pos = JSON.parse(ent.position);
            // ctx.fillText(pos.x, 500, 36 + 50 + 30 * i)
            ctx.fillRect((pos.x / 100) * state.ctx.canvas.width, (pos.y / 100) * state.ctx.canvas.height, 10, 10)

    }
  }


}

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

function drawRoom(points: number[][], s: number) {
  let ctx = state.ctx

  ctx.fillStyle = "rgba(255,149,0,0.1)"

  ctx.lineWidth = 2
  ctx.lineJoin = "round"


  let minX = Math.min(...points.map(point => point[0]));
  let minY = Math.min(...points.map(point => point[1]));
  let maxX = Math.max(...points.map(point => point[0]));
  let maxY = Math.max(...points.map(point => point[1]));

  let scaleX = ctx.canvas.width / (maxX - minX);
  let scaleY = ctx.canvas.height / (maxY - minY);
  let scale = Math.min(scaleX, scaleY);

  let centerX = (maxX + minX) / 2;
  let centerY = (maxY + minY) / 2;
  ctx.save();
  ctx.translate(ctx.canvas.width / 2, ctx.canvas.height / 2);
  // ctx.rotate(315 * Math.PI / 180);

  ctx.beginPath();
  ctx.moveTo((points[0][0] - centerX) * scale, (points[0][1] - centerY) * scale);

  for (let i = 1; i < points.length; i++) {
    ctx.lineTo((points[i][0] - centerX) * scale, (points[i][1] - centerY) * scale);
  }
  ctx.closePath()
  ctx.fill()
  ctx.stroke()
  ctx.restore();
}

function drawGrid() {
  let ctx = state.ctx

  ctx.strokeStyle = "rgba(255,255,255,0.25)"

  ctx.lineWidth = 2

  ctx.lineJoin = "round"
  for (let j = 0; j < 24; j++) {
    for (let i = 0; i < 24; i++) {
      let dx = (state.width / 24)
      let dy = (state.height / 24)
      ctx.beginPath()
      ctx.moveTo(dx / 2 + dx * i, dx / 2 + dy * j)
      ctx.lineTo(dx / 2 + dx * i, dx / 2 + dy * j + 1)
      ctx.closePath()
      ctx.stroke()
    }
  }

}

function redraw() {
  state.ctx.clearRect(0, 0, state.width, state.height)
  drawGrid()
  drawZone()
  state.animationFrame = requestAnimationFrame(redraw)
}

function setupCanvas() {
  const chart = document.getElementById(`zone-canvas`) as HTMLCanvasElement
  if (!chart) return;

  const ctx = chart.getContext('2d')
  if (!ctx) return
  state.ctx = ctx
  let scale = 1.5
  ctx.scale(scale, scale)


  chart.width = chart.clientWidth * scale
  chart.height = chart.clientHeight * scale

  state.width = chart.width
  state.height = chart.height
  // ctx.translate(0, 0)
  // ctx.translate(0.5, 0.5);
  ctx.clearRect(0, 0, state.width, state.height)

  redraw()


}

</script>

<template>
  <Toolbar class="mb-1" icon="􀐘" title="Zones">
    <ToolbarButton :active="false" :disabled="false" class=" px-3"
                   style="height: 1.4rem"
                   text="Select"
                   @click="() => {}"></ToolbarButton>
    <div class="button-sep"></div>
      <ToolbarButton :active="false" :disabled="state.selected.length === 0"
                     class="  px-3"
                     style="height: 1.5rem"
                     text="Delete"
                     @click="() => {deleteZone()}"></ToolbarButton>
    <div class="flex-grow-1"></div>

    <div class="button-sep"></div>
    <ToolbarButton :active="false" icon="􀅼" style="height: 1.5rem"
                   text="Zone"
                   to="/terminal/settings/subroutines/zones/create"></ToolbarButton>
  </Toolbar>
  <div class="d-flex gap-1">

    <div v-if="false" class="d-flex justify-content-start py-2 px-0">
      <div class="label-w500 label-o4 label-xxl"><i
          :class="`fa-solid fa-map fa-fw`"></i></div>
      <div class="label-w500 opacity-100 label-xxl px-2">Zones</div>
      <div class="flex-fill"></div>


      <Plot :cols="2" :rows="1" small style="width: 13rem;">


        <Radio :active="false" :fn="() => toggleShowDeleted()"
               :title="state.showDeleted?'Hide Deleted':'Show Deleted'"></Radio>
        <Radio :active="false"
               :fn="() => state.mode === 'create'?setMode('list'):setMode('create')"
               :title="state.mode === 'create'?'Cancel':'New Zone'"></Radio>
      </Plot>

    </div>
    <div v-if="state.loading">
      <div class="element p-2">
        <div class="label-c1 label-o4 d-flex align-content-center gap-1">
          <div>
            <Loader size="sm"></Loader>
          </div>
          <div class="">Loading...</div>
        </div>
      </div>
    </div>
    <div v-else-if="state.mode === 'list'">
      <div v-if="false" class="pane-container">
        <div class="element p-1 ">

          <div v-for="zone in state.zones" :key="zone.id"
               @click="() => selectZone(zone.id)">
            <div
                :class="zone.id === state.selectedZone.id?'subplot':'subplot subplot-inline'"
                class="d-flex justify-content-between align-items-center">
              <div class="d-flex align-items-start">
                <span class="label-c4 label-o2 lh-1 pt-1"
                      style="width: 0.4rem; margin-top: 6px;"><span
                    v-if="zone.pinned">􀎧</span></span>

                <div class="py-1 px-1">
                    <div
                            class="label-c1 text-capitalize label-w600 label-o5 lh-1">
                        {{ zone.name }}
                        <div v-if="true"
                             class="label-c2 label-o4 label-w500 text-accent">􀷙
                        </div>
                        <div v-else
                             class="label-c2 label-o1 label-w500">􀓞
                        </div>
                    </div>

                  <div class="label-c3 label-w600 label-o3">
                    {{ zone.entities.length }} items
                  </div>
                </div>
              </div>
              <div class="label-c1 label-o3 px-1">
                􀆊
              </div>
            </div>
            <div class="h-sep my-1"
                 style="width: 95%; margin-left: 0.5rem"></div>
          </div>

        </div>

        <div class="zone-spec gap-1">

          <div class="element p-1 px-1">
            <div class="d-flex justify-content-between">

              <div class="info-bar">
                <div
                    class="label-c1 lh-1 label-r label-w500 label-o4">Pinned</div>
                <div
                    class="label-c2 label-o2 lh-1">A pinned zone will show up first on any list of
                                                   zones.
                </div>
              </div>
              <div v-if="state.selectedZone">
                <div v-if="!state.selectedZone.pinned"
                     class="subplot subplot-button"
                     @click="(e) => pinZone()">Pin
                </div>
                <div v-else class="subplot subplot-button"
                     @click="(e) => unpinZone()">Unpin
                </div>
              </div>
            </div>
          </div>

          <div class="element p-1 px-1">
            <div class="d-flex justify-content-between mb-1">
              <div class="info-bar">
                <div
                    class="label-c1 lh-1 label-r label-w500 label-o4">Entities</div>
                <div
                    class="label-c2 label-o2">The following entities are included in this zone.</div>
              </div>
              <div>
                <div v-if="state.listMode === 'view'"
                     class="subplot subplot-button"
                     @click="(e) => setListMode('modify')">Modify
                </div>
                <div v-if="state.listMode === 'modify'"
                     class="subplot subplot-button"
                     @click="(e) => setListMode('view')">Done
                </div>
              </div>
            </div>
            <div v-if="state.selectedZone" class="">
              <FixedScroll class="zones  mx-1">
                <div v-for="entity in state.selectedZone.entities"
                     v-if="state.listMode === 'view'"
                     class="subplot p-2 ">
                  <div class="d-flex gap-1">
                    <div>
                      {{ entity.icon }}
                    </div>
                    <div class="text-capitalize">
                      {{ entity.name }}
                    </div>
                  </div>
                  <span class="label-o3">{{ entity.module }}</span>
                </div>
                <div
                    v-for="entity in state.entities"
                    v-if="state.listMode === 'modify'"
                    :class="state.selectedZone.entities.find(en => en.id === entity.id)?'':'subplot-inline'"
                    class="subplot p-2"
                    @click="(e) => toggleEntity(entity)">
                  <div class="d-flex gap-1">
                    <div>
                      {{ entity.icon }}
                    </div>
                    <div class="text-capitalize">
                      {{ entity.name }}
                    </div>
                  </div>
                  <span class="label-o3">{{ entity.module }}</span>
                </div>
              </FixedScroll>

            </div>
          </div>

          <div class="element p-1 px-1">
            <div class="d-flex justify-content-between">
              <div class="info-bar">
                <div class="label-c1 label-w500 label-o4">Visibility</div>
                <div
                    class="label-c2 label-o2 label-w500">Deleted zones will be made invisible, but
                                                         can be recovered later.
                </div>
              </div>
              <div v-if="state.selectedZone">
                <div v-if="!state.selectedZone.deleted"
                     class="subplot subplot-button"
                     @click="(e) => deleteZone()">Hide
                </div>
                <div v-else class="subplot subplot-button"
                     @click="(e) => restoreZone()">Show
                </div>
              </div>

            </div>
          </div>

        </div>
      </div>

    </div>
    <div v-else-if="state.mode === 'create'"
         class="d-flex justify-content-center btn-outline-primary">
        <CreateZone :done="() => state.mode = 'list'">
        </CreateZone>
    </div>
      <div class="zone-grid">
          <div v-for="z in state.zones" :key="z.id"
               :class="state.selected === z.id?'accent-selected':''"
               class="element haptic d-flex gap-1 flex-column  p-2"
               style="grid-column: span 2;"
               @click="() => selectZone(z.id)">
              <div class="d-flex gap-1">
                  <div class="d-flex flex-row justify-content-between w-100">
                      <div v-if=z.name class="label-c2 label-o4 label-w500 label-r text-capitalize"
                           style="">
                          {{ z.name }}
                      </div>
                      <div v-else class="label-c2 label-o3 label-w500 label-r text-capitalize">
                          Unnamed
                      </div>
                      <div v-if="state.selectedZone?.id === z.id"
                           class="label-c2 label-o4 label-w500 text-accent">􀷙
                      </div>
                      <div v-else
                           class="label-c2 label-o1 label-w500">􀓞
                      </div>

                  </div>


              </div>
              <div class="label-c3 label-o3 label-w5002 ">
                  {{ z.entities.map(e => e.icon).join(' ') }}
              </div>
          </div>
      </div>
      <div v-if=false class="element " style="width: 40rem">

          <div class="d-flex justify-content-between">

          </div>
          <canvas id="zone-canvas" style="height: 24rem; width: 100%"></canvas>
      </div>
  </div>
</template>

<style lang="scss" scoped>
div.element {

}

div.haptic.element:active {
  transform: scale(98%) !important;
}

.zone-grid {
  width: 100%;
  display: grid;
  grid-gap: 0.25rem;
  grid-template-columns: repeat(12, minmax(1rem, 1fr));
  grid-template-rows: repeat(5, minmax(2.5rem, 1fr));
}

.info-bar {
  margin-left: 0.25rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.info-bar > div {
  line-height: 0.75rem;
}

.scroll-element {
  height: 12rem;
  padding-right: 0.25rem;
  margin-bottom: 0.25rem;
  overflow: scroll;
}

.subplot-selected {
  border: 1px solid white !important;
}


.subplot {
  max-height: 2rem;
  border: 1px solid transparent;
}

.zone-spec {
  grid-column: 2 / 5;
  grid-row: 1 / 1;
  display: flex;
  flex-direction: column;
}

.pane-container {
  display: grid;
  grid-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(4, 1fr);
}

.pane-2x1 {
  grid-column: 2 / 5;
  grid-row: 1 / 1;
}

.zones {
  display: grid;
  grid-gap: 0.25rem;
  grid-template-rows: 1.8rem;
  grid-template-columns: repeat(4, 1fr);
}

</style>