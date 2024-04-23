<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {useRouter} from "vue-router";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Entity, Task, TaskOption, Zone} from "@/types";
import {TaskType} from "@/types";
import TaskManager from "@/components/task/TaskManager.vue";
import zoneService from "@/services/zoneService";
import type {Remote} from "@/remote";
import Button from "@/components/Button.vue";

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

})

function setOptions() {
  if (!remote) return;

  let entities: TaskOption[] = remote.entities.map(t => {
    return {title: t.alias ? t.alias : t.name.slice(0, 24), description: t.module, value: t.id}
  }) as TaskOption[]

  state.tasks = [
    {
      title: "Name",
      description: "The name of the Macro",
      type: TaskType.String,
      value: "",
      preview: "Unnamed"
    },
    {
      title: "Entities",
      description: "A brief description of the macro's function.",
      type: TaskType.List,
      options: entities,
      value: [],
      preview: "None"
    },

  ]
  state.loaded = true
}


function goBack() {
  router.push("/terminal/settings/subroutines")
}

function finish(tasks: Task[]) {
  const name = tasks.find(t => t.title === "Name");
  if (!name) return;

  const entities = tasks.find(t => t.title === "Entities");
  if (!entities) return;

  zoneService.createZone({
    name: name.value as string,
    entities: remote.entities.filter(e => (entities.value as string[]).includes(e.id)) as Entity[],
  } as Zone).then(res => {
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
  <div class="">
    <div class="nav-grid gap-1 pb-1 w-100">
      <div class="d-flex justify-content-start">
        <Button accent active fill icon="􀆉" text="Back"
                @click="() => goBack()"></Button>
      </div>
      <div class="d-flex justify-content-center">
        <div class="label-w500 label-c1 label-w600 align-self-center">Create Zone</div>
      </div>
      <div class="d-flex justify-content-end gap-1 align-items-center">
        <!--            <div class="v-sep" style="height: 1rem"></div>-->

      </div>
    </div>
    <div class="d-flex justify-content-center">
      <div class="generic-grid ">
        <TaskManager :on-complete="finish" :tasks="state.tasks" :title="`Zone`"></TaskManager>
      </div>
    </div>
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
  width: max(60%, 0%);
}
</style>