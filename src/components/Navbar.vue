<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>


import router from "@/router";
import Element from "udap-ui/components/Element.vue";

const props = defineProps<{
  back?: string
  title?: string
  cb?: () => void
  auto?: boolean
}>()

//style="border-bottom: 1px solid rgba(255,255,255,0.0); padding-bottom: 0.375rem; margin-bottom: 0.375rem"

function back() {
  if (props.auto) {
    router.go(-1)
    return;
  }
  if (props.cb) {
    props.cb()
    return
  }
  router.go(-1)
}

</script>

<template>
  <div class="d-flex justify-content-between align-items-center w-100">
    <Element :cb="back"
             class="d-flex align-items-center justify-content-center gap-0 py-2 "
             foreground
             mutable
             style="height: 2.5rem; width: 25%">
      <div class="sf-icon label-o1" style="margin-left: -8px">􀯶</div>
      <div class="label-c5 label-w500">Back</div>
    </Element>
    <div class="d-flex justify-content-center label-w600 text-capitalize label-c3 label-o6 " style="width: 50%">
      {{
        props.title
      }}
    </div>
    <div class=" justify-content-end" style="width: 25%">
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