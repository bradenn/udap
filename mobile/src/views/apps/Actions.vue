<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import ElementLink from "udap-ui/components/ElementLink.vue";
import type {SelectOption} from "udap-ui/components/Select.vue";

import {onMounted, reactive, watchEffect} from "vue";
import core from "udap-ui/core";
import {Action, Trigger} from "udap-ui/types";


const remote = core.remote();


const state = reactive({
  actions: [] as Action[],
  options: [] as SelectOption[],

  triggers: [] as Trigger[],

  switches: [] as Trigger[][],

  ready: false
})

onMounted(() => {
  updateActions()
  state.ready = true
})

watchEffect(() => {

  updateActions()
  return remote.actions
})


function updateActions() {
  state.triggers = remote.triggers.filter(s => s.name.startsWith('aa0'))
  state.switches = [[], [], [], [], []];
  for (let i = 0; i < state.triggers.length; i++) {
    let trigger = state.triggers[i]
    let id = trigger.name.replace("mqtt-aa0", "")
    let chunks = id.split("-")
    let device = parseInt(chunks[0]) - 1;
    let button = parseInt(chunks[1]);
    state.switches[device].push(trigger)
  }

  for (let i = 0; i < state.switches.length; i++) {
    state.switches[i] = state.switches[i].sort((a: Trigger, b: Trigger) => a.name.localeCompare(b.name))
  }

  state.actions = remote.actions
}


</script>

<template>
  <Element class="scrollable-fixed">
    <List scroll-y>
      <List row>
        <Element class="justify-content-center d-flex label-w600 label-c5 text-success w-100 align-items-center"
                 foreground mutable
                 style="height: 3rem"
                 to="/apps/actions/create"
        >
          Create Action
        </Element>

      </List>

      <ElementLink v-for="action in state.actions" :title="action.name" :to="`/apps/actions/${action.id}`"
                   icon="􀋦"></ElementLink>


    </List>
  </Element>

  <!--  <List>-->
  <!--    <div class="d-flex gap-1 flex-column">-->
  <!--      <Element class="d-flex justify-content-end">-->
  <!--        <Navbar title="Actions" back="/apps">-->
  <!--        </Navbar>-->
  <!--      </Element>-->
  <!--    </div>-->


  <!--  </List>-->
  <!--  <Element>-->

  <!--    <div v-if="false" class="h-100 mt-1">-->
  <!--      <Element>-->
  <!--        <List>-->

  <!--          <div v-for="device in state.switches">-->

  <!--            <ElementHeader :title="`Switch ${state.switches.indexOf(device) + 1}`"></ElementHeader>-->
  <!--            <List>-->
  <!--              <Element foreground mutable v-for="button in device" v-if="device.length > 0">-->
  <!--                Button {{ button.name.slice(button.name.length - 1) }}-->
  <!--              </Element>-->
  <!--              <Element v-else foreground>No Actions</Element>-->
  <!--            </List>-->
  <!--          </div>-->
  <!--        </List>-->

  <!--      </Element>-->
  <!--    </div>-->


  <!--  </Element>-->
</template>


<style lang="scss" scoped>

</style>