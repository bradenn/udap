<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import core from "@/core";
import type {Entity, Task, TaskOption, Zone} from "@/types";
import {TaskType} from "@/types";
import {onMounted, reactive, watchEffect} from "vue";
import TaskManager from "@/components/task/TaskManager.vue";
import entityService from "@/services/entityService";

onMounted(() => {
  findEntity()
  if (state.loaded) {
    sync()
    setOptions()
  }
})

const router = core.router()
const remote = core.remote()
const state = reactive({
  entity: {} as Entity,
  zones: [] as Zone[],
  tasks: [] as Task[],
  loaded: false,
})

function findEntity() {
  const id = router.currentRoute.value.params["entityId"];
  if (!id) return;
  const sr = remote.entities.find(s => s.id === id)
  if (!sr) return
  state.entity = sr
  state.loaded = true
  sync()
  setOptions()
}

watchEffect(() => {

  if (state.loaded) {
    sync()
    setOptions()
  }
  return remote
})

function setOptions() {
  if (!remote) return;
  let iconSet: string[] = ["􁎶", "􁏀", "􁎿", "􁌡", "􀛭", "􁏁", "􁎦", "􁏙", "􁌥", "􁒜", "􁎾", "􁒪", "􁌜", "􁓼", "􁌦", "􁒰", "􀏞", "􀛮", "􀙫", "􀍉", "􀎲", "􀲰", "􀢹", "􀧘", "􀪯", "􁁋", "􀎚", "􁃗", "􀍽", "􀬗", "􀵔"]
  let icons: TaskOption[] = iconSet.map(t => {
    return {title: t, description: t, value: t}
  }) as TaskOption[]
  state.tasks = [
    {
      title: "Alias",
      description: "A nickname for the entity",
      type: TaskType.String,
      value: state.entity.alias,
      preview: state.entity.alias
    },
    {
      title: "Icon",
      description: "Pick an icon to easily identify the entity",
      type: TaskType.Icon,
      options: icons,
      value: state.entity.icon,
      preview: state.entity.icon
    },
    {
      title: "Position",
      description: "Define a worldspace location for the entity",
      type: TaskType.String,
      value: state.entity.position,
      preview: state.entity.position
    },
  ]
  state.loaded = true
}

function sync() {

}

const notify = core.notify();


function finish(tasks: Task[]) {
  const alias = tasks.find(t => t.title === "Alias");
  if (!alias) return;

  const icon = tasks.find(t => t.title === "Icon");
  if (!icon) return;

  const position = tasks.find(t => t.title === "Position");
  if (!position) return;

  let en = state.entity as Entity
  en.alias = alias.value
  en.icon = icon.value
  en.position = position.value


  entityService.update(en).then((res: any) => {
    notify.success("Entity", "Entity has been updated.")
    done()
  }).catch((err: any) => {
    notify.fail("Entity", "Entity has NOT been updated.", err)
    console.log(err)
  })
}

function save() {
  try {
    finish(state.tasks)
  } catch (err) {
    console.trace(err)
    console.log(err)
  }
}

function done() {
  router.back()
}


</script>

<template>

  <div class="context-grid">
    <div v-if="state.loaded" class="context-pane  d-flex flex-column" style=""
         @click.stop>
      <div class="nav-grid gap-1 pb-1 w-100 px-2">
        <div class="d-flex justify-content-start">
          <div class="label-w500 label-c1 text-accent"
               @click="() => {done()}">􀆉 Back
          </div>

        </div>
        <div class="d-flex justify-content-center">
          <div
              class="label-w500 label-c1 label-w600 align-self-center">Edit Entity
          </div>
        </div>
        <div class="d-flex justify-content-end">
          <div class="label-w500 label-c1 text-accent"
               @click="(e) => {save()}">Save
          </div>
        </div>
      </div>

      <TaskManager :on-complete="finish" :tasks="state.tasks"
                   :title="`Edit`"></TaskManager>
    </div>
  </div>

</template>

<style scoped>
.context-pane {
  grid-column: 3 / span 6;
  grid-row: 2 / span 5;
}

.nav-grid {
  width: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(3, 1fr);
}

.context-grid {
  width: 100%;
  height: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(8, 1fr);
  grid-template-columns: repeat(10, 1fr);
}

.ctx {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 2;
  /*background-color: rgba(0, 0, 0, 0.1);*/
  backdrop-filter: blur(18px) brightness(90%);
}

.preview {
  width: 9rem;
}

.grid-element {
  width: 100%;
  height: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(3, 1fr);
  grid-template-columns: repeat(4, 1fr);
}

.generic-grid {
  display: flex;
  justify-content: center;
  grid-column-gap: 0.25rem;
}
</style>