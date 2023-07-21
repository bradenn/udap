<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import TableItem from "@/components/TableItem.vue";

const props = defineProps<{
    title: string
    icon: string
}>()

</script>

<template>
    <TableItem :title="props.title" icon="􀅴">
        <div class="lh-1 label-c4 label-o4 label-w500">
            <slot></slot>
        </div>
    </TableItem>
</template>

<style lang="scss" scoped>
</style>