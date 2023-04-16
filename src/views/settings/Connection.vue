<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Table from "@/components/Table.vue";
import TableText from "@/components/TableText.vue";
import core from "@/core";
import {onBeforeMount, reactive} from "vue";
import type {Endpoint} from "@/types";


const remote = core.remote()

const state = reactive({
    endpoint: {} as Endpoint,
    storage: {}
})

let token = localStorage.getItem("token")

interface TokenData {
    id: string
}

function parseJwt(token: string): TokenData {
    let base64Url = token.split('.')[1];
    let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    let jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload) as TokenData;
}

onBeforeMount(() => {
    let tokenMeta = parseJwt(token)
    let ep = remote.endpoints.find(e => e.id === tokenMeta.id)
    if (!ep) return
    state.endpoint = ep
    state.storage = localStorage.getItem("preferences")
})


</script>

<template>
    <div class="d-flex flex-column gap-3">
        <Table title="Client">
            <TableText icon="" title="Status">{{ remote.connected ? "Connected" : "Disconnected" }}</TableText>
            <TableText icon="" title="Endpoint Identifier">{{ state.endpoint.name }}</TableText>
            <TableText icon="" title="Permissions">
                <div v-if="state.endpoint.type === 'terminal'">
                    Inexorable
                </div>
                <div v-else>

                </div>
            </TableText>
            <TableText icon="" title="Local Storage">{{ state.storage }}</TableText>
        </Table>

        <Table title="Node">
            <TableText icon="" title="Hostname">{{ remote.metadata.system.hostname }}</TableText>
            <TableText icon="" title="IPv4">{{ remote.metadata.system.ipv4 }}</TableText>
            <TableText icon="" title="Cores">{{ remote.metadata.system.cores }}</TableText>
            <TableText icon="" title="Environment">{{ remote.metadata.system.environment }}</TableText>
            <TableText icon="" title="Version">{{ remote.metadata.system.version }}</TableText>
            <TableText icon="" title="Compiler">{{ remote.metadata.system.go }}</TableText>
        </Table>
    </div>
</template>

<style scoped>

</style>