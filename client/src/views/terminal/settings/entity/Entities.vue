<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Remote} from "@/types";
import Loader from "@/components/Loader.vue";
import Plot from "@/components/plot/Plot.vue";
import PaneMenu from "@/components/pane/PaneMenu.vue";
import PaneMenuItem from "@/components/pane/PaneMenuItem.vue";
import PaneDialogue from "@/components/pane/PaneDialogue.vue";
import PanePopup from "@/components/pane/PanePopup.vue";
import entityService from "@/services/entityService";

let remote = inject('remote') as Remote
let preferences = inject('preferences')

let state = reactive({
  modules: {} as any,
  attributes: {} as Attribute[],
  loading: true,
  selectedEntity: {} as Entity,
  moduleEntities: [] as Entity[],
  selectedAttributes: [] as Attribute[],
  selectedModule: "",
  configOption: "",
  error: ""
})

const icons = ["􀛮", "􁎿", "􁌥", "􁏀", "􀙫", "􀍉", "􀎲", "􀲰", "􀢹", "􀧘", "􀪯", "􀥔", "􁁋", "􀎚", "􁃗", "􀍽", "􀬗", "􀵔"]

onMounted(() => {
  state.loading = true
  handleUpdates(remote)
  selectModule(Object.keys(state.modules)[0])
  state.loading = false
})

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
  state.modules = groupBy<Entity>(remote.entities, 'module') as Entity[]
  state.attributes = groupBy<Attribute>(remote.attributes, 'entity') as Attribute[]

  return remote
}

// groupBy creates several arrays of elements based on the value of a key
function groupBy<T>(xs: T[], key: string): T[] {
  return xs.reduce(function (rv: any, x: any): T {
    (rv[x[key]] = rv[x[key]] || []).push(x);
    return rv;
  }, {});
}

function selectModule(module: string) {
  closePopup()
  state.selectedModule = module
  let entities = remote.entities.filter(e => e.module === module)
  if (!entities) return
  state.moduleEntities = entities

}

function selectEntity(entity: string) {
  closePopup()
  const find = remote.entities.find(e => e.id === entity)
  if (!find) return
  state.selectedEntity = find
  const attrs = remote.attributes.filter(e => e.entity === entity)
  if (!attrs) return
  state.selectedAttributes = attrs
}

function setConfig(conf: string) {
  state.configOption = conf
}

function closePopup(): void {
  state.configOption = ''
}

function isOpen(key: string): boolean {
  return state.configOption === key
}

function setIcon() {
  entityService.setIcon(state.selectedEntity.id, state.selectedEntity.icon).then(res => {
    state.error = JSON.stringify(res)
  }).catch(err => {
    state.error = JSON.stringify(err)
  })
}

function parsePosition(pos: string) {
  interface pos {
    x: number
    y: number
  }

  let pso = JSON.parse(pos) as pos

  return `X: ${pso.x}, Y: ${pso.y}`
}
</script>

<template>
  <div>
    <div v-if="!state.loading" class="entity-grid ">

      <PaneMenu
          :alt="`${Object.keys(state.modules).length} ${Object.keys(state.modules).length === 1?'module':'modules'}`"
          title="Module">

        <PaneMenuItem v-for="(entities, module) in state.modules"
                      :active="state.selectedModule === `${module}`"
                      :fn="() => selectModule(`${module}`)"
                      :subtext="`${entities.length} ${entities.length === 1?'entity':'entities'}`"
                      :title="`${module}`"></PaneMenuItem>
      </PaneMenu>
      <PaneMenu :alt="`${state.moduleEntities.length} ${state.moduleEntities.length === 1?'entity':'entities'}`"
                title="Entities">

        <PaneMenuItem v-for="entity in state.moduleEntities" :active="state.selectedEntity.id === entity.id"
                      :fn="() => selectEntity(entity.id)"
                      :subtext="entity.type"
                      :title="entity.name"></PaneMenuItem>
      </PaneMenu>
      <div v-if="state.selectedEntity.id" class="pane-fill gap-1">
        <PaneDialogue :alt="state.selectedEntity.alias" subtext="An alias will replace the default name of an entity."
                      title="Alias">
          <div class="subplot subplot-button">
            Create Alias
          </div>
        </PaneDialogue>
        <PaneDialogue :alt="state.selectedEntity.icon"
                      subtext="The icon assigned by default can be modified to better represent the entity's purpose"
                      title="Icon">
          <div class="subplot subplot-button" @click="(e) => setConfig('icon')">
            Configure
          </div>
        </PaneDialogue>
        <PanePopup v-if="isOpen('icon')" :apply="setIcon"
                   :close="closePopup" title="Select an Icon">
          <div class="icon-grid">
            <div v-for="icon in icons"
                 :class="`${state.selectedEntity.icon === icon?'subplot-active':'subplot-inline'}`"
                 class="subplot subplot-button"
                 @click="state.selectedEntity.icon = icon">
              {{ icon }}
            </div>
          </div>
        </PanePopup>
        <PaneDialogue :alt="state.selectedEntity.neural"
                      subtext="Allow UDAP's neural network to determine the appropriate state for this light at any given time"
                      title="Neural">
          <div class="subplot subplot-button">
            Enable
          </div>
        </PaneDialogue>
        <PaneDialogue v-if="state.selectedEntity?.position"
                      :alt="state.selectedEntity?.position !== '{}'?parsePosition(state.selectedEntity?.position):''"
                      subtext="Assign coordinates to the approximate location of this entity within a zone."
                      title="Position">
          <div class="subplot subplot-button" @click="(e) => setConfig('position')">
            Configure
          </div>
        </PaneDialogue>
        <PanePopup v-if="isOpen('position')" :close="closePopup"
                   title="Position">
        </PanePopup>
        <PaneDialogue v-if="state.selectedEntity?.position"
                      :alt="state.selectedEntity?.position !== '{}'?parsePosition(state.selectedEntity?.position):''"
                      subtext="Assign default states for attributes in the event of a power loss"
                      title="Defaults">
          <div class="subplot subplot-button" @click="(e) => setConfig('defaults')">
            Configure
          </div>
        </PaneDialogue>
        <PanePopup v-if="isOpen('defaults')" :close="closePopup"
                   title="Defaults">
          <div v-for="attribute in state.selectedAttributes" :key="attribute.id">
            {{ attribute.key }}
            {{ attribute.value }}
          </div>
        </PanePopup>
        <PaneDialogue v-if="state.selectedEntity?.position"
                      alt="Peak: 15 W, Linear"
                      subtext="Assign coordinates to the approximate location of this entity within a zone."
                      title="Energy Usage">
          <div class="subplot subplot-button" @click="(e) => setConfig('power')">
            Configure
          </div>
        </PaneDialogue>
        <PanePopup v-if="isOpen('power')" :close="closePopup"
                   title="Power Usage">
        </PanePopup>
      </div>

      <div v-for="(entities, module) in state.modules" v-if="false">
        <div v-if="entities" class="">
          <Plot :cols="2" :rows="2" :title="`${module}`">
            <div v-for="entity in entities.slice(0,4)" class="subplot">
              <div class="d-flex justify-content-start align-items-center flex-row px-1 w-100">
                <div class="label-w500 label-o3 label-c1">{{ entity.icon }}&nbsp;</div>
                <div class="label-w500 label-c1">{{ entity.name }}</div>
                <div class="flex-grow-1"></div>
                <div class="label-w500 label-c2 label-o2"><i class="fa-solid fa-gear"></i></div>
              </div>
            </div>
          </Plot>
        </div>

      </div>


    </div>
    <div v-else class="element p-2">
      <div class="label-c1 label-o4 d-flex align-content-center gap-1">
        <div>
          <Loader size="sm"></Loader>
        </div>
        <div class="">Loading...</div>
      </div>
    </div>
  </div>
</template>

<style scoped>


.pane-fill {
  position: relative;
  grid-column: 3 / 6;
  grid-row: 1 / 1;
  display: flex;
  flex-direction: column;
}

.icon-grid {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(8, 1fr);
}

.entity-grid {
  width: 100%;
  height: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(5, 1fr);
}

</style>