<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import core from "@/core";
import type {Macro, Task, TaskOption, Zone} from "@/types";
import {TaskType} from "@/types";
import {onMounted, reactive, watchEffect} from "vue";
import TaskManager from "@/components/task/TaskManager.vue";
import macroService from "@/services/macroService";

interface EditMacroProps {
  macro: Macro,
  done?: () => void
}

const props = defineProps<EditMacroProps>()

onMounted(() => {
  sync()
  setOptions()
})

const router = core.router()
const remote = core.remote()
const state = reactive({
  zone: {} as Zone,
  zones: [] as Zone[],
  tasks: [] as Task[],
  loaded: false,
})

watchEffect(() => {
  sync()
  setOptions()
  return remote.macros
})

function setOptions() {
  if (!remote) return;

  let zones: TaskOption[] = remote.zones.map(t => {
    return {title: t.name, description: t.entities.map(e => e.name).join(", "), value: t.id}
  }) as TaskOption[]

  state.tasks = [
    {
      title: "Name",
      description: "What should the macro be named",
      type: TaskType.String,
      value: props.macro.name,
      preview: props.macro.name
    },
    {
      title: "Description",
      description: "Briefly describe the macro's function",
      type: TaskType.String,
      value: props.macro.description,
      preview: props.macro.description.length + " chars"
    },
    {
      title: "Zone",
      description: "What zone should this macro target",
      type: TaskType.Radio,
      options: zones,
      value: state.zone.id,
      preview: state.zone.name
    },
    {
      title: "Attribute",
      description: "Which attribute key should be selected",
      type: TaskType.Radio,
      options: [
        {
          title: "Power",
          value: 'on',
          description: "Change the power state"
        }, {
          title: "Brightness",
          value: 'dim',
          description: "Change the brightness"
        }, {
          title: "Color",
          value: 'hue',
          description: "Change the color hue"
        }, {
          title: "Color Temperature",
          value: 'cct',
          description: "Change the color temperature"
        },
      ],
      value: props.macro.type,
      preview: props.macro.type
    },
    {
      title: "Value",
      description: "What value should be requested",
      type: TaskType.String,
      value: props.macro.value,
      preview: props.macro.value
    },
  ]
  state.loaded = true
}

function sync() {
  const zone = remote.zones.find(z => z.id === props.macro.zone)
  if (!zone) return
  state.zone = zone
  state.zones = remote.zones

}

function finish(tasks: Task[]) {
  const name = tasks.find(t => t.title === "Name");
  if (!name) return;

  const desc = tasks.find(t => t.title === "Description");
  if (!desc) return;

  const zone = tasks.find(t => t.title === "Zone");
  if (!zone) return;

  const key = tasks.find(t => t.title === "Attribute");
  if (!key) return;

  const value = tasks.find(t => t.title === "Value");
  if (!value) return;

  macroService.updateMacro({
    name: name.value as string,
    zone: zone.value,
    type: key.value,
    description: desc.value,
    id: props.macro.id,
    value: value.value,
  } as Macro).then(res => {
    if (props.done) {
      props.done()
    }
  }).catch(err => {
    console.log(err)
  })
}

function save() {
  finish(state.tasks)
}


</script>

<template>
  <div class="ctx ">
    <div class="context-grid">
      <div v-if="state.loaded" class="context-pane  d-flex flex-column" style="" @click.stop>
        <div class="nav-grid gap-1 pb-1 w-100 px-2">
          <div class="d-flex justify-content-start">
            <div class="label-w500 label-c1 text-accent" @click="() => {if(props.done) props.done()}">􀆉 Back</div>

          </div>
          <div class="d-flex justify-content-center">
            <div class="label-w500 label-c1 label-w600 align-self-center">Edit Macro</div>
          </div>
          <div class="d-flex justify-content-end">
            <div class="label-w500 label-c1 text-accent" @click="save">Save</div>
          </div>
        </div>

        <TaskManager :on-complete="finish" :tasks="state.tasks" :title="`Edit`"></TaskManager>
      </div>
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