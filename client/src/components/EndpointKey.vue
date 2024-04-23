<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>


let props = defineProps<{
    value: string
}>()


</script>

<template>
    <div class="d-flex justify-content-center gap-1 align-items-center">
        <div class="surface char">
            {{ props.value.slice(0, 4) }}
            &nbsp;
            {{ props.value.slice(4, 8) }}
        </div>
    </div>
</template>

<style lang="scss" scoped>
.separator {
  display: flex;
  justify-content: center;
  height: 100%;
  padding: 0.0625rem 0.125rem;
  font-weight: 600;
  border-radius: 0.25rem;
  color: rgba(255, 255, 255, 0.3) !important;
  font-size: 0.8rem;
}

.char {
  display: flex;
  justify-content: center;
  height: 100%;
  padding: 0.0625rem 0.33rem;
  border-radius: 0.25rem;
  color: rgba(255, 255, 255, 0.6) !important;
  font-size: 0.8rem;

}
</style>