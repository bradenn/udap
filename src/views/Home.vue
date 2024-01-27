<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import core from "udap-ui/core";
import {onMounted, reactive, watchEffect} from "vue";
import type {Entity as EntityType, Zone} from "udap-ui/types";
import Entity from "@/components/Entity.vue";
import ThermostatV2 from "@/components/ThermostatV2.vue";
import Element from "udap-ui/components/Element.vue";
import ElementLink from "udap-ui/components/ElementLink.vue";
import List from "udap-ui/components/List.vue";
import {PreferencesRemote} from "udap-ui/persistent";
import ActionToggle from "udap-ui/components/ActionToggle.vue";


const remote = core.remote()

const preferences = core.preferences() as PreferencesRemote

const state = reactive({
  entities: [] as EntityType[],
  zones: [] as Zone[],
  thermostat: {} as EntityType,
})


onMounted(() => {

  state.entities = remote.entities
  state.zones = remote.zones.filter(z => z.pinned)
})

watchEffect(() => {
  state.zones = remote.zones.filter(z => z.pinned)
  let e = remote.entities.find(e => e.name === "thermostat")
  if (e) state.thermostat = e
  return remote.entities
})


</script>


<template>
  <div class="d-flex flex-column gap-3 justify-content-between h-100">

    <div v-if="preferences.landscape" class="d-flex flex-column gap-2 w-50">
      <div>
        <div class="label-c5 label-w700 label-o5 px-2">Climate</div>
        <div class="w-100">
          <ThermostatV2 :entity="state.thermostat"></ThermostatV2>

        </div>
      </div>
      <div class="flex-column d-flex gap-2">
        <div v-for="zone in state.zones">
          <div>
            <div class="label-c5 label-w700 label-o5 px-2">{{ zone.name }}</div>
            <div class="home-grid-xl">
              <Entity v-for="e in zone.entities" :key="e.id" :entity="e"></Entity>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="d-flex flex-column gap-2">
      <div>
        <div class="label-c5 label-w700 label-o5 px-2">Climate</div>
        <div class="w-100">
          <ThermostatV2 :entity="state.thermostat"></ThermostatV2>
          <!--          <br>-->
          <!--          <ThermostatV3 :entity="state.thermostat"></ThermostatV3>-->
        </div>
      </div>
      <List>
        <div v-for="zone in state.zones">
          <div class="label-c5 label-w700 label-o5 px-2">{{ zone.name }}</div>
          <div class="home-grid">
            <Entity v-for="e in zone.entities" :key="e.id" :entity="e"></Entity>
          </div>
        </div>
      </List>
      <div class="label-c5 label-w700 label-o5 px-2">Actions</div>
      <div class="home-grid">
        <ActionToggle v-for="a in preferences.home" v-if="preferences.home" :id="a"></ActionToggle>
        <div v-else style="grid-column: 1 / span 2;">
          <Element>
            <ElementLink icon="􁖇" title="Select Actions" to="/apps/actions"></ElementLink>
          </Element>
        </div>
      </div>
      <div v-if="false">
        <div class="label-c5 label-w700 label-o5 px-2 mb-0 pb-0 lh-0">Notifications</div>
        <Element>
          <ElementLink icon="􁖇" title="Homework Assignment" to="/extra/homework"></ElementLink>
        </Element>
      </div>
    </div>
  </div>

</template>

<style scoped>
.home-grid-xl {
  display: grid;
  gap: 0.25rem;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  grid-template-rows: repeat(1, minmax(0, 1fr));
}

.home-grid {
  display: grid;
  gap: 0.25rem;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  grid-template-rows: repeat(1, minmax(0, 1fr));
}
</style>