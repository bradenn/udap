<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import core from "@/core";
import {computed, reactive} from "vue";
import Navbar from "@/components/Navbar.vue";

const router = core.router()

const state = reactive({
    currentName: ""
})


router.afterEach((to, from, failure) => {
    if (!to.name) return
    state.currentName = <string>to.name
})

let transitionName = computed(() => {
    if (router.currentRoute.value.path === '/home/settings') {
        return 'slide-reverse'
    }
    return 'slide'
})

</script>

<template>
    <div v-if="transitionName !== 'slide-reverse'" class="d-flex">
        <Navbar v-if="router.currentRoute.value.name" :title="state.currentName"
                back="/home/settings">
            Hello
        </Navbar>
    </div>
    <router-view v-slot="{ Component }">
        <transition :name=transitionName mode="out-in">
            <component :is="Component"/>
        </transition>
    </router-view>
</template>

<style scoped>


@keyframes slide-in {
    from {
        transform: translateX(10%);
        filter: blur(12px);
        opacity: 0;
    }
    to {
        transform: translateX(0) scale(1);
        filter: blur(0px);
        opacity: 1;
    }
}

@keyframes slide-out {
    from {
        transform: translateX(0) scale(1);
        filter: blur(0px);
        opacity: 1;
    }
    to {
        transform: translateX(-10%);
        filter: blur(12px);
        opacity: 0;
    }
}

.slide-enter-active {
    animation: slide-in 100ms ease-out;
}

.slide-leave-active {
    animation: slide-out 100ms ease-out;
}

/*.slide-enter-to {*/
/*    transform: translateX(0);*/
/*    opacity: 1;*/
/*}*/

/*.slide-leave-to {*/
/*    transform: translateX(10%);*/
/*    opacity: 0;*/
/*}*/
</style>