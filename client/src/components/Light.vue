<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watch} from "vue";
import AttributeComponent from "@/components/entity/Attribute.vue"
import type {Attribute, Entity} from "@/types"

// Establish a local reactive state
let state = reactive<{
  active: boolean,
  showMenu: boolean,
  activeColor: string,
  powerAttribute: Attribute,
  attributes: Attribute[]
}>({
  active: false,
  showMenu: false,
  activeColor: "rgba(255,255,255,1)",
  powerAttribute: {} as Attribute,
  attributes: []
})

// Define the prop for this entity
let props = defineProps<{
  entity: Entity
}>()

// Inject the remote manifest
const remote: any = inject('remote')

// When the view loads, force the local state to update
onMounted(() => {
  updateLight(remote.attributes)
})

// Ripple any attribute changes down to the local reactive state
watch(remote.attributes, (newAttributes: Attribute[]) => {
  updateLight(newAttributes)
})

// Compare the defined order to sort the lights
function compareOrder(a: any, b: any): number {
  return a.order - b.order
}

// Update the reactive model for the light
function updateLight(attributes: Attribute[]): void {
  // Define the attributes for the light
  state.attributes = attributes.filter((a: Attribute) => a.entity === props.entity.id).sort(compareOrder)
  // Get the current power state of the light
  let current = attributes.find((a: Attribute) => a.entity === props.entity.id && a.key === 'on')
  // Assign the power state to a local attribute
  state.powerAttribute = current as Attribute
  // Register weather the light is on or off
  state.active = current ? current.value === 'true' : false;
}

// Toggle the state of the context menu
function toggleMenu(): void {
  state.showMenu = !state.showMenu
}
</script>

<template>
  <div>
    <div :class="state.active?'active':''" class="entity-small" @click="toggleMenu">
      <div class="entity-header mb-2 ">
        <div class="label-o5">
          {{ props.entity.icon || '􀛮' }}
        </div>
        <div class="label-c1 label-w400 label-o4 px-2">
          {{ props.entity.name }}
        </div>
      </div>
      <div>

      </div>
      <div class="fill"></div>
      <div class="label-xxs label-o3 label-w500 px-2">
        <div v-if="state.active">ON</div>
        <div v-else>OFF</div>
      </div>
    </div>
    <div v-if="state.showMenu" class="context" @click="toggleMenu"></div>
    <div v-if="state.showMenu" @click="toggleMenu">
      <div class="entity-context top d-flex">
        <div class="d-flex flex-column gap px-3 top ">
          <div class="d-flex justify-content-start align-items-end align-content-end" v-on:click.stop>
            <div>
              <span :style="`text-shadow: 0 0 8px ${state.activeColor};`"
                    class="label-md label-w600 label-o3">{{ props.entity.icon }}</span>
              <span class="label-md label-w600 label-o6 px-2">{{ props.entity.name }}</span>
            </div>
            <div class="fill"></div>
            <div class="h-bar">
              <AttributeComponent :attribute="state.powerAttribute" :entity-id="props.entity.id"
                                  small></AttributeComponent>
            </div>
          </div>

          <div class="context-container gap v-bar">
            <div
                v-for="attribute in state.attributes.filter((attribute: Attribute) => attribute.key !== 'on')"
                :key="attribute.id">
              <AttributeComponent :key="attribute.id" :attribute="attribute" :entity-id="props.entity.id" primitive
                                  small></AttributeComponent>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>