<script>

import Dock from "../Dock.vue";
import VueQrcode from 'vue-qrcode'

export default {
  name: "Wifi",
  components: {Dock, VueQrcode},
  data() {
    return {
      options: {
        data: "hello",
        cellSize: 9,
      }
    }

  },
  beforeCreate() {

  },
  created() {

  },
  mounted() {

  },
  beforeDestroy() {

  },
  methods: {},
}
</script>

<template>
  <div class="d-flex flex-column gap justify-content-center align-items-center align-content-center h-100">
    <div class="element p-3 pt-0">
      <div class="label-ys label-o3 py-1 pt-2">Scan to connect to Wi-Fi network '<span class="label-w600">udap</span>'
      </div>
      <vue-qrcode :color="{ dark: '#ffffff50', light: '#00000000' }" :margin=0 :quality=1 :scale=8 class=""
                  errorCorrectionLevel='H' type="image/png" value="WIFI:T:WPA;S:udap;P:Zaq1xsw2;"></vue-qrcode>
    </div>
  </div>
</template>

<style scoped>

</style>
