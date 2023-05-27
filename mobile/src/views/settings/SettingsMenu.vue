<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import TableLink from "@/components/TableLink.vue";
import Table from "@/components/Table.vue";
</script>

<template>
    <div class="d-flex flex-column gap-3">

        <Table title="System">
            <TableLink icon="􀣋" title="General" to="/home/settings/general"></TableLink>
            <TableLink icon="􀉣" title="Connection" to="/home/settings/connection"></TableLink>
        </Table>


        <Table title="Appearance">
            <TableLink icon="􀏅" title="Background" to="/home/settings/background"></TableLink>
            <TableLink icon="􀟗" title="Color Theme" to="/home/settings/color"></TableLink>
        </Table>
    </div>
</template>

<style scoped>
</style>