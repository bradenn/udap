<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import core from "@/core";
import type {Macro, SubRoutine, Task, TaskOption, Trigger, Zone} from "@/types";
import {TaskType} from "@/types";
import {onMounted, reactive, watchEffect} from "vue";
import TaskManager from "@/components/task/TaskManager.vue";
import subroutineService from "@/services/subroutineService";

interface EditMacroProps {
  subroutine: SubRoutine
  done?: () => void
}

const props = defineProps<EditMacroProps>()

onMounted(() => {
  sync()
  setOptions()
})

const router = core.router()
const remote = core.remote()
const notifications = core.notifications()
const state = reactive({
  zone: {} as Zone,
  triggers: [] as Trigger[],
  macros: [] as Macro[],
  tasks: [] as Task[],
  loaded: false,
})

watchEffect(() => {
  sync()
  setOptions()

})


function setOptions() {
  if (!remote) return;
  let iconSet: string[] = ["􁎶", "􁏀", "􁎿", "􁌡", "􀛭", "􁏁", "􁎦", "􁏙", "􁌥", "􁒜", "􁎾", "􁒪", "􁌜", "􁓼", "􁌦", "􁒰", "􀏞"]
  let triggers: TaskOption[] = remote.triggers.map(t => {
    return {title: t.name, description: t.description, value: t.id}
  }) as TaskOption[]

  let macros: TaskOption[] = remote.macros.map(t => {
    return {title: t.name, description: t.description, value: t.id}
  }) as TaskOption[]

  let icons: TaskOption[] = iconSet.map(t => {
    return {title: t, description: t, value: t}
  }) as TaskOption[]

  let trigger = remote.triggers.find(t => t.id === props.subroutine.triggerId)
  if (!trigger) return;

  state.tasks = [
    {
      title: "Name",
      description: "What should the subroutine be named?",
      type: TaskType.String,
      value: props.subroutine.description,
      preview: props.subroutine.description
    },
    {
      title: "Trigger",
      description: "What trigger should evoke this subroutine?",
      type: TaskType.Radio,
      options: triggers,
      value: props.subroutine.triggerId,
      preview: trigger.name
    },
    {
      title: "Icon",
      description: "Which icon?",
      type: TaskType.Icon,
      options: icons,
      value: props.subroutine.icon,
      preview: `${props.subroutine.icon}`
    },
    {
      title: "Macros",
      description: "Which macros should be run when this subroutine is triggered?",
      type: TaskType.List,
      options: macros,
      value: props.subroutine.macros.map(m => m.id),
      preview: `${props.subroutine.macros.length} items`
    },
    {
      title: "Revert",
      description: "Should this event return to a default state after a duration?",
      type: TaskType.Number,
      value: props.subroutine.revertAfter,
      preview: calcDuration(props.subroutine.revertAfter)
    }
  ]
  state.loaded = true
}

function calcDuration(min: number) {
  return `${Math.floor(min / 60)}h ${min % 60}m `
}

function sync() {
  const macros = remote.macros
  if (!macros) return
  state.macros = macros
  const triggers = remote.triggers
  if (!triggers) return
  state.triggers = triggers

}

function finish(tasks: Task[]) {
  const name = tasks.find(t => t.title === "Name");
  if (!name) return;

  const trigger = tasks.find(t => t.title === "Trigger");
  if (!trigger) return;

  const macros = tasks.find(t => t.title === "Macros");
  if (!macros) return;

  const icon = tasks.find(t => t.title === "Icon");
  if (!icon) return;

  const revert = tasks.find(t => t.title === "Revert");
  if (!revert) return;

  let macroValues = macros.value as string[]
  let macroTypes = macroValues.map(m => remote.macros.find(mc => mc.id === m)) as Macro[]

  subroutineService.updateSubroutine({
    description: name.value as string,
    icon: icon.value as string,
    macros: macroTypes as Macro[],
    revertAfter: revert.value as number,
    triggerId: trigger.value as string,
    id: props.subroutine.id,
  } as SubRoutine).then(_ => {
    if (props.done) props.done()
    notifications.show("Subroutine", "Update succeeded!!!", 1, 1000 * 8)
  }).catch(err => {
    notifications.show("Subroutine", "HOLY FUCK", 2, 1000 * 8)
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
            <div class="label-w500 label-c1 label-w600 align-self-center">{{ props.subroutine.description }}
            </div>
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
.accent-selected {
  outline: 2px solid rgba(255, 149, 0, 0.4) !important;
}

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