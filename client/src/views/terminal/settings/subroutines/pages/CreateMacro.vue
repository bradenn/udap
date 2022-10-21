<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {useRouter} from "vue-router";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Macro, Remote, Task, TaskOption} from "@/types";
import {TaskType} from "@/types";
import TaskManager from "@/components/task/TaskManager.vue";
import macroService from "@/services/macroService";

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

  let zones: TaskOption[] = remote.zones.filter(z => !z.deleted).map(t => {
    return {title: t.name, description: t.entities.map(e => e.name).join(", "), value: t.id}
  }) as TaskOption[]

  state.tasks = [
    {
      title: "Name",
      description: "The name of the Macro",
      type: TaskType.String,
      value: "",
      preview: ""
    },
    {
      title: "Description",
      description: "A brief description of the macro's function.",
      type: TaskType.String,
      value: "",
      preview: ""
    },
    {
      title: "Zone",
      description: "What devices should this macro interact with?",
      type: TaskType.Radio,
      options: zones,
      value: "",
      preview: "unset"
    },
    {
      title: "Attribute Target",
      description: "Which attributes should be targeted within the zone?",
      type: TaskType.Radio,
      options: [{
        title: "Brightness",
        description: "Change the brightness of a light",
        value: "dim"
      }, {
        title: "Color",
        description: "Change the color of a light",
        value: "hue"
      }, {
        title: "Color Temperature",
        description: "Change the color of a light",
        value: "cct"
      }, {
        title: "Power",
        description: "Change the power state of a device",
        value: "on"
      }],
      value: "",
      preview: ""
    },
    {
      title: "Value",
      description: "What should the attribute values be set to?",
      type: TaskType.String,
      value: "",
      preview: ""
    }
  ]
  state.loaded = true
}

function goBack() {
  router.push("/terminal/settings/subroutines")
}

function finish(tasks: Task[]) {
  const name = tasks.find(t => t.title === "Name");
  if (!name) return;

  const description = tasks.find(t => t.title === "Description");
  if (!description) return;

  const zone = tasks.find(t => t.title === "Zone");
  if (!zone) return;

  const attributeTarget = tasks.find(t => t.title === "Attribute Target");
  if (!attributeTarget) return;

  const value = tasks.find(t => t.title === "Value");
  if (!value) return;

  macroService.createMacro({
    name: name.value as string,
    description: description.value as string,
    zone: zone.value as string,
    type: attributeTarget.value as string,
    value: value.value as string,
  } as Macro).then(res => {
    console.log(res)
    goBack()
  }).catch(err => {
    console.log(err)
  })

}


function createSubroutine() {


}

</script>

<template>
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