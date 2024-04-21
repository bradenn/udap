<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import ElementSpectral from "udap-ui/components/ElementSpectral.vue";

import type {SelectOption} from "udap-ui/components/Select.vue";
import Select from "udap-ui/components/Select.vue";
import actionService from "udap-ui/services/actionService"
import {onMounted, reactive, watchEffect} from "vue";
import core from "udap-ui/core";
import {Action, Entity, Trigger} from "udap-ui/types";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import ElementTextBox from "udap-ui/components/ElementTextBox.vue";
import ElementPair from "udap-ui/components/ElementPair.vue";
import Navbar from "@/components/Navbar.vue";


const remote = core.remote();

const state = reactive({

  action: {
    attribute: "spectral"
  } as Action,
  step: 0,
  actionType: "" as string,
  triggerOptions: [] as SelectOption[],
  entityOptions: [] as SelectOption[],
  attributeOptions: [] as SelectOption[],
  ready: false
})

onMounted(() => {
  updateOptions()

  state.ready = true
})

watchEffect(() => {
  updateOptions()

})

function updateOptions() {
  state.triggerOptions = buildTriggerOptions(remote.triggers.filter(e => e.name.includes("uhidv2-aa")));
  state.entityOptions = buildEntityOptions(remote.entities.filter(e => (e.type === "spectrum" || e.type === "switch") && e.module != "govee"));
  state.attributeOptions = buildAttributeOptions([]);
}

function updateTrigger(option: SelectOption[]) {
  state.action.triggerId = option[0].id
}

function updateEntity(options: SelectOption[]) {
  state.action.entities = options.map(e => e.id)
}

function updateAttribute(options: SelectOption[]) {
  state.action.attribute = options[0].id
}


function buildEntityOptions(entities: Entity[]): SelectOption[] {
  let options: SelectOption[] = []
  entities.sort((a, b) => {
    return (a.alias ? a.alias : a.name).localeCompare((b.alias ? b.alias : b.name))
  })
  for (let entity of entities) {
    let option = {} as SelectOption
    option.name = entity.alias ? entity.alias : entity.name
    option.description = entity.module
    option.id = entity.id
    option.category = entity.module
    option.icon = entity.icon
    options.push(option)
  }

  return options
}

const attributeOptions = [
  {
    name: "Spectrum",
    id: "spectral",
    description: "Control all light features",
    icon: "􀑇"
  } as SelectOption,

]

function buildAttributeOptions(attributes: string[]): SelectOption[] {
  let options: SelectOption[] = []
  attributes.sort((a, b) => {
    return a.localeCompare(b)
  })
  options.push({
    name: "Spectrum",
    id: "spectral",
    description: "Control all light features",
    icon: "􀑇"
  } as SelectOption)

  // options.push({
  //   name: "Brightness",
  //   id: "dim",
  //   description: "Control light brightness",
  //   icon: "􁐄"
  // } as SelectOption)
  //
  //
  // options.push({
  //   name: "Color Temperature",
  //   id: "cct",
  //   description: "Control light color temp",
  //   icon: "􀡊"
  // } as SelectOption)

  options.push({
    name: "Power",
    id: "on",
    description: "Control power state",
    icon: "􁏮"
  } as SelectOption)

  // options.push({
  //   name: "Hue",
  //   id: "hue",
  //   description: "Control light color",
  //   icon: "􀟗"
  // } as SelectOption)


  return options
}

function buildTriggerOptions(triggers: Trigger[]): SelectOption[] {
  let options: SelectOption[] = []
  triggers.sort((a, b) => {
    return a.name.localeCompare(b.name)
  })
  for (let trigger of triggers) {
    let option = {} as SelectOption
    option.name = trigger.name
    option.description = trigger.description
    option.id = trigger.id
    option.category = trigger.type
    option.icon = "􀋦"
    options.push(option)
  }

  return options
}

function updateSelection(options: SelectOption[]) {
  state.actionType = options[0].id
}

function updateAttributeRequest(spectral: string) {
  state.action.request = spectral
}

function togglePower() {
  if (state.action.request == 'false') {
    state.action.request = 'true'
    state.action.attribute = "on"
  } else {
    state.action.request = 'false'
    state.action.attribute = "on"
  }

}


function proceed() {
  state.step += 1;
  if (state.step > 4) {
    createAction()
  }
}

const router = core.router();

function createAction() {
  actionService.create(state.action).then(() => {
    router.push("/apps/actions")
  }).catch(err => {
    alert(err)
  })
}

function changeText(data: string) {
  state.action.name = data
}

</script>

<template>
  <div class="d-flex flex-column gap-1 scrollable-fixed">

    <Element class="d-flex flex-shrink-0">
      <Navbar back="/apps/actions" title="Create Action">
        <Element :cb="proceed" class="justify-content-center d-flex label-w600" foreground
                 mutable>
          {{ state.step < 3 ? "Next" : "Create" }}
        </Element>
      </Navbar>
    </Element>


    <Select v-if="state.step == 0" :options="state.triggerOptions"
            :selections="updateTrigger"
            single
            title="Select a trigger"></Select>

    <div v-else-if="state.step == 1">
      <Select v-if=state.ready :options="state.entityOptions" :selections="updateEntity" many prompt="pick many"
              title="Select entities"></Select>
    </div>
    <div v-else-if="state.step == 2">
      <Select v-if=state.ready :options="state.attributeOptions" :selections="updateSelection" prompt="pick one" single
              title="Select an attribute"></Select>
    </div>
    <div v-else-if="state.step == 3">
      <List>
        <Element>
          <ElementHeader title="Settings"></ElementHeader>
          <ElementSpectral v-if="state.actionType == 'spectral'" :change="updateAttributeRequest"
                           :value="'on;cct:-1;70'"></ElementSpectral>

          <ElementPair v-else-if="state.actionType=='on'" :cb="togglePower"
                       :value="state.action.request == 'true'?'ON':'OFF'"
                       icon="􀆨"
                       title="Power"></ElementPair>
        </Element>

      </List>
    </div>

    <div v-else-if="state.step == 4">
      <ElementTextBox :change="changeText"></ElementTextBox>
    </div>


  </div>
</template>

<style lang="scss" scoped>

</style>