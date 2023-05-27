<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import TableItem from "@/components/TableItem.vue";

const props = defineProps<{
    title: string
    icon: string
    to: string
}>()

</script>

<template>
    <TableItem :icon="props.icon" :title="props.title" :to="props.to">
        <div class="label-o1 label-w600 label-c4 sf-icon end">􀯻</div>
    </TableItem>
</template>

<style lang="scss" scoped>
</style>