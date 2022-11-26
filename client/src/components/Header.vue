<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

interface Props {
    icon?: string,
    name: string,
    done?: () => void,
    save?: () => void,

}

const props = defineProps<Props>()

</script>

<template>
    <div class="nav-grid gap-1 pb-1 w-100 ">
        <div class="d-flex justify-content-start">
            <div class="label-w500 label-c1 text-accent element subplot py-1 px-3"
                 @click="() => {if(props.done) props.done()}">􀆉 Back</div>

        </div>
        <div class="d-flex justify-content-center">
            <div class="label-w500 label-c1 label-w600 align-self-center gap-1 d-flex ">
                <span v-if="props.icon" class="label-o3">{{ props.icon }}</span>
                {{ props.name }}
            </div>
        </div>
        <div class="d-flex justify-content-end">
            <div class="label-w500 label-c1 text-accent element subplot py-1 px-3" @click="save">Save</div>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.nav-grid {
  width: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(3, 1fr);
}


</style>
