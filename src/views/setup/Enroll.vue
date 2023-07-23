<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>


import {onMounted, reactive} from "vue";
import {Preference} from "@/preferences";
import {PreferenceTypes} from "udap-ui/types";
import axios from "axios";

import Element from "udap-ui/components/Element.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import List from "udap-ui/components/List.vue";


const state = reactive({
  name: "",
  token: "",
  coolDown: 0,
  message: "",
  delta: 0
});

onMounted(() => {
})

function authenticate() {

  // Generated endpoint registration url using the security code and controller address
  let url = `https://api.udap.app/endpoints/register/${state.token}`
  // Make the request to the controller app
  axios.get(url).then(res => {
    // Set the token in localStorage
    new Preference(PreferenceTypes.Token).set(res.data.token)
    new Preference(PreferenceTypes.Name).set(state.name.slice(0, 1).toUpperCase() + state.name.toLowerCase().slice(1))
    // Redirect the user to the authenticated portal
    window.location.href = "/"

  }).catch(err => {
    // alert("Denied..." + err)
    state.message = "Invalid access token"
    console.log(err)
  })
}

</script>

<template>
  <div class="d-flex flex-column gap-1  justify-content-center" style="height: 30vh">


    <Element>
      <ElementHeader class="pt-0 mb-2" title="Access Token">
        <template v-slot:description>
          <div class="text-danger label-c5 lh-1" v-text="state.message"></div>
        </template>
      </ElementHeader>
      <List row>
        <input id="cypher" v-model="state.token" autocapitalize="off" autocomplete="off"
               class="subplot w-100"
               placeholder="--------" type="text">
        <Element :cb="() => authenticate()" accent class="d-flex justify-content-center align-items-center py-1"
                 foreground style="width: 8rem">
          Login
        </Element>
      </List>
    </Element>

  </div>
</template>

<style scoped>
.form-control {
  font-size: 1.2rem;
  font-weight: 500;
}
</style>