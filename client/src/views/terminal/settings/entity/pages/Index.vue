<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity} from "@/types";
import type {Remote} from "@/remote";
import MenuSection from "@/components/menu/MenuSection.vue";
import EntityView from "@/components/Entity.vue";
import FixedScroll from "@/components/scroll/FixedScroll.vue";
import {useRouter} from "vue-router";
import Toolbar from "@/components/toolbar/Toolbar.vue";
import type {RouteMeta} from "@/composables/routeMeta";
import {useRouteMeta} from "@/composables/routeMeta";
import ToolbarButton from "@/components/ToolbarButton.vue";

let remote = inject('remote') as Remote
let preferences = inject('preferences')


let sortModes = new Map<string, string>([["module", "Module"], ["alphabetical", "Alphabetical"], ["recent", "Recent"]])

let state = reactive({
  modules: {} as any,
  mode: "module" as string,
  entities: {} as Entity[],
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
  state.loading = false
})

watchEffect(() => handleUpdates(remote))

function sortEntities(a: Entity, b: Entity): number {
  if (state.mode === 'recent' || state.mode === 'module') {
    if (new Date(a.updated).valueOf() < new Date(b.updated).valueOf()) {
      return 1
    } else if (new Date(a.updated).valueOf() > new Date(b.updated).valueOf()) {
      return -1
    }
    return 0
  } else if (state.mode === 'alphabetical') {
    if (a.name.toUpperCase() < b.name.toUpperCase()) {
      return -1
    } else if (a.name.toUpperCase() > b.name.toUpperCase()) {
      return 1
    }
    return 0
  }

  return 0
}

function setMode(value: string) {
  loadEntities()
  state.mode = value
}

function loadEntities() {
  state.entities = remote.entities.sort(sortEntities)
}

function handleUpdates(remote: Remote) {
  state.modules = groupBy<Entity>(remote.entities, 'module') as Entity[]
  loadEntities()
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

let meta: RouteMeta = useRouteMeta()


</script>

<template>
  <div class="h-100 d-flex flex-column gap-1">
    <Toolbar :icon="meta?.icon" class="flex-shrink-0"
             title="Entities">
      <div class="flex-fill"></div>

      <div class="label-c1 label-o3 label-w500 label-r">Sort:</div>
      <ToolbarButton v-for="mode in sortModes" :accent="state.mode === mode[0]"
                     :text="mode[1]"
                     @click="() => {setMode(mode[0])}"></ToolbarButton>
    </Toolbar>
    <FixedScroll v-if="state.mode === 'module'"
                 style="overflow-y: scroll !important; height: 100%; max-height: 100% !important;">
      <div v-for="(entities, module) in state.modules"
           class="d-flex flex-column gap-1">
        <MenuSection :title="`${module}`">
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
    <FixedScroll
        v-else-if="state.mode === 'alphabetical' || state.mode === 'recent'"
        style="overflow-y: scroll !important; height: 100%; max-height: 100% !important;">
      <div
          class="entity-grid">
        <EntityView
            v-for="entity in state.entities"
            :entity="entity"
            :noselect="true"
            style="min-width: 6rem"
            @click="() => goToEntity(entity.id)"></EntityView>

      </div>
    </FixedScroll>
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