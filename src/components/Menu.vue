<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import core from "@/core";
import {onBeforeMount, reactive} from "vue";
import AppLink from "@/components/AppLink.vue";

const props = defineProps<{
    name?: number
}>()

const state = reactive({
    meta: {} as any,
    menu: false
})

onBeforeMount(() => {
    state.meta = fetchJWT()
})

const remote = core.remote()

function fetchJWT(): any {
    let token = localStorage.getItem("token")
    if (!token) {
        return ""
    }

    if (token === "unset") {
        return ""
    }

    return parseJwt(token)
}

function parseJwt(token: string): string {
    let base64Url = token.split('.')[1];
    let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    let jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
}

function toggleMenu() {
    state.menu = !state.menu
}

</script>

<template>
    <div class="navbar-dropdown d-flex align-content-center align-items-center justify-content-between element"
         style="z-index: 10 !important;"
         @click="() => toggleMenu()">
        <div class="d-flex align-items-center gap-2 align-items-center">
            <div class="label-c5 label-o5 label-w600 lh-1 px-2">Menu</div>
        </div>
        <div class="label-c2 label-o2" style="margin-top: -8px;">
            ⌵
        </div>
    </div>
    <div v-if="state.menu" class="menu" @click="() => toggleMenu()">
        <div class="blur"></div>
        <div class="context-menu">
            <div class="d-flex gap-1">
                <AppLink icon="􀎟" name="Home" to="/home/dashboard"></AppLink>
                <AppLink icon="􀍟" name="Settings" to="/home/settings"></AppLink>
                <!--                <AppLink name="Home" icon="􀎟" to="/home/dashboard"></AppLink>-->
            </div>
        </div>
    </div>
</template>

<style scoped>
.blur {
    -webkit-backdrop-filter: blur(4px);
    width: calc(100% + 1rem);
    height: calc(100% + 1rem);
    position: absolute;
    top: -0.5rem;
    left: -0.5rem;
    z-index: 5 !important;
}

.context-menu {
    /*outline: 1px solid red;*/
    position: relative;
    top: 3.5rem;
    left: 0;
    height: 6rem;
    z-index: 2000;
    animation: blur-in ease-in-out 200ms;
}

@keyframes blur-in {
    0% {
        opacity: 0;
    }
    100% {
        opacity: 1;
    }
}

.menu {
    position: absolute;
    width: calc(100% - 1rem);
    height: 100%;
}

.navbar-dropdown {
    /*border: 1px solid rgba(255, 255, 255, 0.1);*/
    /*border-radius: 0.5rem;*/
    padding: 0.25rem 0.70rem;
    width: 8rem;
}

.navbar-dropdown:active {

}
</style>