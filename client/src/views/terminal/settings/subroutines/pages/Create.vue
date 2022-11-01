<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {useRouter} from "vue-router";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Macro, Remote, SubRoutine, Task, TaskOption} from "@/types";
import {TaskType} from "@/types";
import TaskManager from "@/components/task/TaskManager.vue";
import subroutineService from "@/services/subroutineService";
import Page from "@/components/Page.vue";

const router = useRouter()

const remote = inject("remote") as Remote

interface Props {
  done?: () => {},
}

const props = defineProps<Props>()

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
      preview: "unset"
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
      preview: "0 items"
    },
    {
      title: "Revert",
      description: "Should this event return to a default state after a duration?",
      type: TaskType.Number,
      value: 0,
      preview: "0h 0m"
    }
  ]
  state.loaded = true
}

function goBack() {
  if (props.done) {
    props.done()
  }
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

  const revert = tasks.find(t => t.title === "Revert");
  if (!revert) return;

  let macroValues = macros.value as string[]
  let macroTypes = macroValues.map(m => remote.macros.find(mc => mc.id === m)) as Macro[]

  subroutineService.createSubroutine({
    description: name.value as string,
    macros: macroTypes as Macro[],
    revertAfter: revert.value as number,
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
  <div class="ctx " @click="goBack()">
    <div class="context-grid">
      <div class="d-flex justify-content-center h-100 w-100"
           style="grid-column: 4 / span 6 !important; grid-row: 2 / span 10 !important;">

        <div class="generic-grid h-100" @click.stop>
          <div>
            <Page :done="() => {goBack()}" :save="() => {}" title="New Subroutine"></Page>

            <TaskManager :on-complete="finish" :tasks="state.tasks" :title="`Macro`" class="h-100"></TaskManager>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>


.context-grid {
  width: 100%;
  height: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(12, 1fr);
  grid-template-columns: repeat(12, 1fr);
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

.generic-grid {
  width: 100%;
}

</style>