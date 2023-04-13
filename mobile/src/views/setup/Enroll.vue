<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>


import {reactive} from "vue";
import {Preference} from "@/preferences";
import {PreferenceTypes} from "@/types";
import axios from "axios";

const state = reactive({
    name: "",
    token: ""
});

function authenticate() {
    // Generated endpoint registration url using the security code and controller address
    let url = `https://api.udap.app/endpoints/register/${state.token}`
    // Make the request to the controller app
    axios.get(url).then(res => {
        // Set the token in localStorage
        new Preference(PreferenceTypes.Token).set(res.data.token)
        new Preference(PreferenceTypes.Name).set(state.name.slice(0, 1).toUpperCase() + state.name.toLowerCase().slice(1))
        // Redirect the user to the authenticated portal
        window.location.href = "/mobile/home"
    }).catch(err => {
        alert("Denied..." + err)
        console.log(err)
    })
}

</script>

<template>
    <div class="d-flex flex-column gap-3">
        <h4>Authentication</h4>
        <div class="d-flex flex-column gap-2 w-100">
            <div class="d-flex flex-column">
                <input id="name" v-model="state.name" class="form-control" placeholder="First Name" type="text">
            </div>
            <input id="cypher" v-model="state.token" class="form-control" placeholder="Access Token" type="text">
        </div>
        <div class="element d-flex justify-content-center" @click="authenticate">Authenticate</div>
    </div>
</template>

<style scoped>

</style>