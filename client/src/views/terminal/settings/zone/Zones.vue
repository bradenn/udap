<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Entity, Remote, Zone} from "@/types";
import Loader from "@/components/Loader.vue";
import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import CreateZone from "@/views/terminal/settings/zone/CreateZone.vue";
import axios from "axios";
import FixedScroll from "@/components/scroll/FixedScroll.vue";

let remote = inject("remote") as Remote
let preferences = inject('preferences')

let state = reactive({
  zones: {} as Zone[],
  entities: [] as Entity[],
  deleted: {} as Zone[],
  selected: "",
  selectedZone: {} as Zone,
  showDeleted: false,
  loading: true,
  mode: "list",
  listMode: "view",
  model: "",
})


onMounted(() => {
  state.loading = true
  handleUpdates()
  state.selectedZone = state.zones[0]
  state.loading = false
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
  state.zones = remote.zones.filter(z => state.showDeleted ? true : !z.deleted).sort((a, b) => sortZones(a, b))
  state.entities = remote.entities

}

function setMode(mode: string) {
  state.mode = mode
}

function selectZone(id: string) {
  state.selected = id
  let target = state.zones.find(z => z.id === id)
  if (!target) return
  state.selectedZone = target
}

function toggleShowDeleted() {
  state.showDeleted = !state.showDeleted
}

function setListMode(value: string) {
  state.listMode = value
}

function deleteZone() {
  let local = state.selectedZone
  if (!local) return
  axios.post(`http://10.0.1.2:3020/zones/${local.id}/delete`)
      .then((r) => {
        state.selectedZone.deleted = true
      })
      .catch((r) => {
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

</script>

<template>
  <div>
    <div v-if="false" class="d-flex justify-content-start py-2 px-0">
      <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-map fa-fw`"></i></div>
      <div class="label-w500 opacity-100 label-xxl px-2">Zones</div>
      <div class="flex-fill"></div>


      <Plot :cols="2" :rows="1" small style="width: 13rem;">


        <Radio :active="false" :fn="() => toggleShowDeleted()"
               :title="state.showDeleted?'Hide Deleted':'Show Deleted'"></Radio>
        <Radio :active="false" :fn="() => state.mode === 'create'?setMode('list'):setMode('create')"
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
      <div class="pane-container">
        <div class="element p-1 ">

          <div v-for="zone in state.zones" :key="zone.id" @click="() => selectZone(zone.id)">
            <div :class="zone.id === state.selectedZone.id?'subplot':'subplot subplot-inline'"
                 class="d-flex justify-content-between align-items-center">
              <div class="d-flex align-items-start">
                <span class="label-c4 label-o2 lh-1 pt-1" style="width: 0.4rem; margin-top: 6px;"><span
                    v-if="zone.pinned">􀎧</span></span>
                <div class="py-1 px-1">
                  <div class="label-c1 text-capitalize label-w600 label-o5 lh-1">{{ zone.name }}</div>

                  <div class="label-c3 label-w600 label-o3">{{ zone.entities.length }} items</div>
                </div>
              </div>
              <div class="label-c1 label-o3 px-1">
                􀆊
              </div>
            </div>
            <div class="h-sep my-1" style="width: 95%; margin-left: 0.5rem"></div>
          </div>

        </div>

        <div class="zone-spec gap-1">

          <div class="element p-1 px-1">
            <div class="d-flex justify-content-between">

              <div class="info-bar">
                <div class="label-c1 lh-1 label-r label-w500 label-o4">Pinned</div>
                <div class="label-c2 label-o2 lh-1">A pinned zone will show up first on any list of zones.</div>
              </div>
              <div v-if="state.selectedZone">
                <div v-if="!state.selectedZone.pinned" class="subplot subplot-button"
                     @click="(e) => pinZone()">Pin</div>
                <div v-else class="subplot subplot-button"
                     @click="(e) => unpinZone()">Unpin</div>
              </div>
            </div>
          </div>

          <div class="element p-1 px-1">
            <div class="d-flex justify-content-between mb-1">
              <div class="info-bar">
                <div class="label-c1 lh-1 label-r label-w500 label-o4">Entities</div>
                <div class="label-c2 label-o2">The following entities are included in this zone.</div>
              </div>
              <div>
                <div v-if="state.listMode === 'view'" class="subplot subplot-button"
                     @click="(e) => setListMode('modify')">Modify</div>
                <div v-if="state.listMode === 'modify'" class="subplot subplot-button"
                     @click="(e) => setListMode('view')">Done</div>
              </div>
            </div>
            <div v-if="state.selectedZone" class="">
              <FixedScroll class="zones  mx-1">
                <div v-for="entity in state.selectedZone.entities" v-if="state.listMode === 'view'"
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
                    class="label-c2 label-o2 label-w500">Deleted zones will be made invisible, but can be recovered later.</div>
              </div>
              <div v-if="state.selectedZone">
                <div v-if="!state.selectedZone.deleted" class="subplot subplot-button"
                     @click="(e) => deleteZone()">Hide</div>
                <div v-else class="subplot subplot-button"
                     @click="(e) => restoreZone()">Show</div>
              </div>

            </div>
          </div>

        </div>
      </div>

    </div>
    <div v-else-if="state.mode === 'create'" class="d-flex justify-content-center btn-outline-primary">
      <CreateZone :done="() => state.mode = 'list'">
      </CreateZone>
    </div>
    <div v-else>


    </div>
  </div>
</template>

<style lang="scss" scoped>

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