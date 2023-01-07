<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity} from "@/types";
import Plot from "@/components/plot/Plot.vue";
import PaneMenu from "@/components/pane/PaneMenu.vue";
import PaneMenuItem from "@/components/pane/PaneMenuItem.vue";
import PaneDialogue from "@/components/pane/PaneDialogue.vue";
import PanePopup from "@/components/pane/PanePopup.vue";
import type {Remote} from "@/remote";
import Input from "@/views/Input.vue";
import MenuSection from "@/components/menu/MenuSection.vue";
import EntityView from "@/components/Entity.vue";
import FixedScroll from "@/components/scroll/FixedScroll.vue";
import {useRouter} from "vue-router";

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


const router = useRouter();

function goToEntity(entityId: string) {
  router.push(`/terminal/settings/entities/${entityId}`)
}


</script>

<template>
  <div class="h-100">
    <FixedScroll
        style="overflow-y: scroll !important; height: 100%; max-height: 100% !important;">
      <div v-for="(entities, module) in state.modules"
           class="d-flex flex-column gap-1">
        <MenuSection :title="module">
          <FixedScroll :horizontal="true"
                       style="max-width: 100%; overflow-x: scroll">
            <div class="d-flex flex-row gap-1">
              <EntityView
                  v-for="entity in entities"
                  :entity="entity"
                  :noselect="true"
                  style="min-width: 6rem"
                  @click="() => goToEntity(entity.id)"></EntityView>
            </div>
          </FixedScroll>
        </MenuSection>

      </div>
    </FixedScroll>

    <div v-if="false" class="entity-grid">

      <PaneMenu
          :alt="`${Object.keys(state.modules).length} ${Object.keys(state.modules).length === 1?'module':'modules'}`"
          title="Module">

        <PaneMenuItem v-for="(entities, module) in state.modules"
                      :active="state.selectedModule === `${module}`"
                      :fn="() => selectModule(`${module}`)"
                      :subtext="`${entities.length} ${entities.length === 1?'entity':'entities'}`"
                      :title="`${module}`"></PaneMenuItem>
      </PaneMenu>
      <PaneMenu
          :alt="`${state.moduleEntities.length} ${state.moduleEntities.length === 1?'entity':'entities'}`"
          title="Entities">

        <PaneMenuItem v-for="entity in state.moduleEntities"
                      :active="state.selectedEntity.id === entity.id"
                      :fn="() => selectEntity(entity.id)"
                      :subtext="entity.type"
                      :title="entity.name"></PaneMenuItem>
      </PaneMenu>
      <div v-if="state.selectedEntity.id" class="pane-fill gap-1">
        <PaneDialogue :alt="state.selectedEntity.alias"
                      subtext="An alias will replace the default name of an entity."
                      title="Alias">
          <div class="subplot subplot-button"
               @click="(e) => setConfig('alias')">
            Create Alias
          </div>
        </PaneDialogue>
        <Input v-if="isOpen('alias')" :apply="setAlias" :close="closePopup"
               :value="state.selectedEntity.alias"
               description="Change the alias" name="Alias"
               style="z-index: 1000 !important;"></Input>

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
          <div class="subplot subplot-button"
               @click="(e) => setConfig('position')">
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
          <div class="subplot subplot-button"
               @click="(e) => setConfig('defaults')">
            Configure
          </div>
        </PaneDialogue>
        <PanePopup v-if="isOpen('defaults')" :close="closePopup"
                   title="Defaults">
          <div v-for="attribute in state.selectedAttributes"
               :key="attribute.id">
            {{ attribute.key }}
            {{ attribute.value }}
          </div>
        </PanePopup>
        <PaneDialogue v-if="state.selectedEntity?.position"
                      alt="Peak: 15 W, Linear"
                      subtext="Assign coordinates to the approximate location of this entity within a zone."
                      title="Energy Usage">
          <div class="subplot subplot-button"
               @click="(e) => setConfig('power')">
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
              <div
                  class="d-flex justify-content-start align-items-center flex-row px-1 w-100">
                <div class="label-w500 label-o3 label-c1">{{
                    entity.icon
                  }}&nbsp;
                </div>
                <div class="label-w500 label-c1">{{ entity.name }}</div>
                <div class="flex-grow-1"></div>
                <div class="label-w500 label-c2 label-o2"><i
                    class="fa-solid fa-gear"></i></div>
              </div>
            </div>
          </Plot>
        </div>

      </div>


    </div>
  </div>
</template>

<style lang="scss" scoped>

.entity-grid {
  width: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-columns: repeat(8, 1fr);

}

.layout-grid {
  width: 100%;

  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-columns: repeat(12, 1fr);
  grid-template-rows: repeat(3, 1fr);
}

.layout-sidebar {
  grid-column: 10 / span 3;
}

.layout-body {

  grid-column: 0 / span 11;

  grid-row: 0 / span 2;
}

</style>