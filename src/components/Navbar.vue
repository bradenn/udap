<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
const props = defineProps<{
    back?: string
    title?: string
}>()
</script>

<template>
    <div class="navbar w-100 mb-3">
      <div class="cell justify-content-start">
        <router-link :to="props.back" class="back label-c4">
          <div class="sf-icon label-o1">􀯶</div>
          <div>Back</div>
        </router-link>
      </div>
      <div class="cell justify-content-center label-w600 text-capitalize label-c4 label-o6">{{
          props.title
        }}
      </div>
      <div class="cell justify-content-end">
        <slot></slot>
      </div>
    </div>
</template>

<style scoped>
.back {
    width: 100%;
    /*outline: 1px solid red;*/
    display: flex;
    flex-direction: row;
    text-decoration: none !important;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.6) !important;
}

.cell {
  width: 33%;
  display: flex;
  align-items: center;
}

.navbar {
    display: flex;
    justify-content: space-between;
}

.navbar > div {
  /*outline: 1px solid rebeccapurple;*/
}
</style>