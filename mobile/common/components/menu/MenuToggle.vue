<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {computed} from 'vue'

const props = defineProps(['modelValue'])
const emit = defineEmits(['update:modelValue'])

const value = computed({
  get() {
    return props.modelValue
  },
  set(value) {
    emit('update:modelValue', value)
  }
})

</script>

<template>

</template>

<style lang="scss" scoped>

</style>