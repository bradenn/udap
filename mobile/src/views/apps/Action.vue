<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import ElementToggle from "udap-ui/components/ElementToggle.vue";
import ElementPair from "udap-ui/components/ElementPair.vue";

import type {SelectOption} from "udap-ui/components/Select.vue";
import actionService from "udap-ui/services/actionService"
import {onMounted, reactive, watchEffect} from "vue";
import core from "udap-ui/core";
import {Action, Entity, Trigger} from "udap-ui/types";
import Navbar from "@/components/Navbar.vue";
import ElementSpectral from "udap-ui/components/ElementSpectral.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";


const remote = core.remote();
const router = core.router();
const state = reactive({
  id: "",
  remote: "" as string,
  trigger: {} as Trigger,
  entities: [] as Entity[],
  action: {} as Action,
  onHome: false,
  edited: false,
  ready: false
})

const preferences = core.preferences();

onMounted(() => {
  updateActions(remote.actions)
})

watchEffect(() => {
  return updateActions(remote.actions)
})

watchEffect(() => {
  checkOnHome(state.action?.id || "")
  return preferences.home
})

function checkOnHome(id: string) {
  if (!preferences.home) {
    preferences["home"] = [id]
  }
  state.onHome = !!preferences.home.find(h => h == id)
}

function toggleOnHome(id: string) {
  if (!preferences.home) {
    preferences["home"] = [id]
  }
  if (preferences.home.find(h => h == id)) {
    preferences.home = preferences.home.filter(h => h != id)
    state.onHome = false;
  } else {
    preferences.home.push(id)
    state.onHome = true;
  }
}

function updateActions(actions: Action[]): boolean {
  state.id = <string>router.currentRoute.value.params["actionId"];
  let action = actions.find(a => a.id == state.id)
  if (!action) return false
  if (state.remote === action.request) {
    return false;
  }

  state.action = action as Action
  console.log(state.action.request)

  state.trigger = remote.triggers.find(t => t.id === action?.triggerId) || {} as Trigger
  state.entities = action.entities.map(eid => remote.entities.find(e => e.id == eid) || {} as Entity)

  state.ready = true
  return true
}

function deleteAction() {
  actionService.delete(state.action).then(ok => {
    router.push("/apps/actions")
  }).catch(err => {
    alert(`Could not delete action: ${err}`)
  })
}

function updateParams(value: string) {
  state.action.request = value

  // state.edited = state.remote.request !== state.action.request;

}

function pushChanges() {
  actionService.update(state.action).then(() => {
    console.log("Done")
  }).catch((err) => console.log(err))
}

function trigger(id: string) {
  actionService.trigger(id).then((err) => {
    console.log(err)
  }).catch((err) => {
    console.log(err)
  })
}

function fitName(name: string) {
  return name

}
</script>

<template>
  <div class="h-100">
    <Element class="d-flex justify-content-end">
      <Navbar :title="`${fitName(state.action.name)}`" back="/apps/actions">
        <Element v-if=state.edited :cb="pushChanges"
                 class="justify-content-center d-flex label-w600 label-c5 text-success"
                 foreground
                 mutable>
          Save
        </Element>

        <Element v-else :cb="deleteAction" class="justify-content-center d-flex label-w600 label-c5 text-danger"
                 foreground
                 mutable>
          Delete
        </Element>
      </Navbar>
    </Element>
    <List class="mt-1" scroll-y style="max-height: 82vh">
      <Element>

        <ElementHeader title="Trigger"></ElementHeader>
        <div class="d-flex flex-column gap-0">
          <List row>
            <ElementPair :title="state.trigger.name" icon="􀒘"></ElementPair>
            <Element :cb="() => trigger(state.action.id)"
                     class="d-flex align-items-center justify-content-center flex-grow-0" foreground
                     mutable style="width: 6rem !important;">
              <div class="sf-icon label-o4" style="margin-left: -8px">􀋦</div>
              Run
            </Element>

          </List>
          <div>
            <ElementToggle :cb="() => toggleOnHome(state.action.id)" :selected="state.onHome"
                           alt="Add this action to your home view"
                           class="mt-1"
                           icon="􁋁" title="Shown in home"></ElementToggle>
          </div>
        </div>
        <ElementHeader class="mt-1" title="Action"></ElementHeader>


        <ElementSpectral v-if="state.action.attribute == 'spectrum'" :change="updateParams"
                         :value="state.action.request" class="pb-1"></ElementSpectral>
        {{ state.action }}

      </Element>

      <Element>
        <List>
          <ElementHeader class="mt-1" title="Devices"></ElementHeader>
          <List>
            <ElementPair v-for="entity in state.entities" :icon="entity.icon"
                         :title="entity.alias?entity.alias:entity.name"
                         button
            ></ElementPair>
          </List>
        </List>
      </Element>

    </List>
  </div>
</template>


<style lang="scss" scoped>

</style>