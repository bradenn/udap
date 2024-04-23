<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

const props = defineProps<{
  name?: string
  value: number
  min: number
  max: number
  warning?: number
  danger?: number
  critical?: number
}>()

</script>

<template>

  <div class="d-flex flex-column">
    <div v-if="props.name" class="d-flex justify-content-between px-1">
      <div class="label-c6 label-o6 label-w600">{{ props.name }}</div>
      <div class="label-c6 label-o6 label-w600">
        <slot></slot>
      </div>
    </div>
    <div class="utilization-container">
      <div :style="`width:${((props.value-props.min)/(props.max-props.min))*100}%;`" class="utilization-bar">
      </div>
    </div>
  </div>

</template>

<style lang="scss" scoped>
.utilization-bar {
  height: 100%;
  color: rgba(255, 255, 255, 0.7);
  //box-shadow: inset 0 0 2px 0.5px rgba(255, 255, 255, 0.05) !important;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 1rem;
  min-width: 4px;

}

.utilization-container {
  border: 1px solid rgba(46, 46, 48, 0.4);
  display: flex;
  padding: 0.065rem;
  align-items: center;
  height: 0.5rem;
  width: 100%;
  //background-color: rgba(255, 255, 255, 0.1);
  background-color: rgba(46, 46, 48, 0.4) !important;
  color: rgba(255, 255, 255, 0.9) !important;

  .sf-icon {
    color: rgba(255, 255, 255, 0.9) !important;
  }

  border-radius: 1rem;
}
</style>