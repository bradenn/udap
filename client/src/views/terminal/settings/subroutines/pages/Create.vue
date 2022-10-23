<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {useRouter} from "vue-router";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Macro, Remote, SubRoutine, Task, TaskOption} from "@/types";
import {TaskType} from "@/types";
import TaskManager from "@/components/task/TaskManager.vue";
import subroutineService from "@/services/subroutineService";

const router = useRouter()

const remote = inject("remote") as Remote

onMounted(() => {
  setOptions()
})

const state = reactive({
  tasks: [] as Task[],
  loaded: false,
})

watchEffect(() => {
  setOptions()
  return remote
})

function setOptions() {
  if (!remote) return;

  let triggers: TaskOption[] = remote.triggers.map(t => {
    return {title: t.name, description: t.description, value: t.id}
  }) as TaskOption[]

  let macros: TaskOption[] = remote.macros.map(t => {
    return {title: t.name, description: t.description, value: t.id}
  }) as TaskOption[]

  state.tasks = [
    {
      title: "Name",
      description: "What should the subroutine be named?",
      type: TaskType.String,
      value: "",
      preview: "Subroutine 1"
    },
    {
      title: "Trigger",
      description: "What trigger should evoke this subroutine?",
      type: TaskType.Radio,
      options: triggers,
      value: "",
      preview: "unset"
    },
    {
      title: "Macros",
      description: "Which macros should be run when this subroutine is triggered?",
      type: TaskType.List,
      options: macros,
      value: [],
      preview: "unset"
    }
  ]
  state.loaded = true
}

function goBack() {
  router.push("/terminal/settings/subroutines")
}

function createTrigger() {
  router.push("/terminal/settings/subroutines/trigger")
}


function finish(tasks: Task[]) {
  const name = tasks.find(t => t.title === "Name");
  if (!name) return;

  const trigger = tasks.find(t => t.title === "Trigger");
  if (!trigger) return;

  const macros = tasks.find(t => t.title === "Macros");
  if (!macros) return;

  let macroValues = macros.value as string[]
  let macroTypes = macroValues.map(m => remote.macros.find(mc => mc.id === m)) as Macro[]

  subroutineService.createSubroutine({
    description: name.value as string,
    macros: macroTypes as Macro[],
    triggerId: trigger.value as string
  } as SubRoutine).then(res => {
    goBack()
  }).catch(err => {
    console.log(err)
  })

}


function createSubroutine() {


}

</script>

<template>
  <div class="d-flex align-items-start label-o4 gap-1 pb-2">
    <div class="label-w500 label-c1 text-accent" @click="goBack">􀆉 Back</div>
  </div>
  <div class="generic-grid">

    <TaskManager :on-complete="finish" :tasks="state.tasks" :title="`Macro`"></TaskManager>
  </div>
</template>

<style scoped>

.generic-grid > div {
  width: 18rem;
}

.generic-grid {
  display: flex;
  justify-content: center;
  grid-column-gap: 0.25rem;
}
</style>