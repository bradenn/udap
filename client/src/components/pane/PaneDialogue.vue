<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Plot {
  title?: string
  subtext?: string
  alt?: string
}

let props = defineProps<Plot>()
</script>

<template>
  <div class="element p-1 px-1">
    <div class="d-flex justify-content-between">
      <div class="info-bar">
        <div class="d-flex align-items-end align-content-end">
          <div class="label-c1 label-w500 label-o4 lh-1">{{ props.title }}</div>
          <div v-if="props.alt" class="label-c2 label-w400 label-o2 mx-1 lh-1">{{ props.alt }}</div>
        </div>
        <div
            class="label-c2 label-o2 label-w500">{{ props.subtext }}
        </div>
      </div>
      <slot></slot>

    </div>
  </div>
</template>


<style lang="scss" scoped>

.info-bar {
  margin-left: 0.25rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 1.65rem;
}

.info-bar > div {
  line-height: 0.75rem;
}

</style>