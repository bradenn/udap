<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import core from "@/core";
import {onBeforeMount, reactive} from "vue";

const props = defineProps<{
    name?: number
}>()

const state = reactive({
    meta: {} as any
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


</script>

<template>
    <div class="navbar-dropdown d-flex align-content-center align-items-center justify-content-between" @click="">
        <div class="d-flex align-items-center gap-2 align-items-center">
            <div class="label-c5 label-o6 label-w500 lh-1"></div>
        </div>
        <div class="label-c2 label-o2" style="margin-top: -8px;">
            ⌵
        </div>

    </div>
</template>

<style scoped>
.navbar-dropdown {
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 0.5rem;
    padding: 0.25rem 0.70rem;
    width: 8rem;
}

.navbar-dropdown:active {

}
</style>