<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import Element from "@/components/Element.vue";

const props = defineProps<{
    title?: string
}>()
</script>

<template>
  <div class="outer-table">
    <div v-if="props.title" class="label-c6 label-w700 label-o3 lh-1 px-3 mb-2 ">{{
        props.title
      }}
    </div>
    <Element class="table">
      <slot></slot>
    </Element>
  </div>
</template>

<style scoped>
.outer-table {
  overflow-x: scroll;
  max-height: 100vh;
  height: 100%;
}

.table {
  /*overflow: scroll !important;*/
  height: 100%;
  padding: 0.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;

}
</style>