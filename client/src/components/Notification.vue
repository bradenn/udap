<!-- Copyright (c) 2022 Braden Nicholson -->


<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import Toast from "@/components/Toast.vue";
import type {Notify, NotifyBody} from "@/notifications";

interface Loader {
    size?: string
}

let props = defineProps<Loader>()

let notify: Notify = inject("notify") as Notify


let state = reactive({
    notify: {} as NotifyBody
})

onMounted(() => {
    updateNotification()
})


watchEffect(() => {
    updateNotification()
    return notify.current()
})

function updateNotification() {
    state.notify = notify.current()
}


</script>

<template>

    <div v-if="notify.isActive()" class="toast-stack ">
        <Toast :key="state.notify.uuid"
               :index="1"
               :message="state.notify.message"
               :severity="state.notify.severity"
               :time="state.notify.duration"
               :title="state.notify.name"></Toast>

    </div>
</template>

<style lang="scss" scoped>

</style>