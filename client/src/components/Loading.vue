<script>

export default {
  components: {},
  data() {
    return {}
  },
  watch: {},
  props: {
    large: Boolean,
    medium: Boolean
  },
  created() {

  },
  methods: {},
}


</script>


<template>
  <div v-if="medium" class="lds-spinner">
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
  </div>
  <div v-else :class="`${large?'lg':''} `" class="sm-ic">
    <div class="lds-grid">
      <div></div>
      <div></div>
      <div></div>
      <div></div>
      <div></div>
      <div></div>
      <div></div>
      <div></div>
      <div></div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.sm-ic {
  zoom: 0.18;
  font-size: 0.7rem;
  line-height: 0.7rem;
}

.sm-ic.lg {
  position: absolute;
  top: 4rem;
  left: 4rem;
  zoom: 1;
  font-size: 4rem;
  line-height: 4rem;
  transition: opacity 1s ease-in-out;
}

.sm-ic.md {
  zoom: 0.24;
  font-size: 1.4rem;
  line-height: 1.4rem;
}

@keyframes lds-grid-done {
  0% {
    opacity: 0.5;
  }
  100% {
    opacity: 0;
    transform: scale(0);
  }
}

.lds-grid {
  display: inline-block;
  position: relative;
  width: 80px;
  height: 80px;

}

.lds-grid div {
  position: absolute;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.5);
  animation: lds-grid 1.2s linear infinite;
}

.lds-grid div:nth-child(1) {
  top: 8px;
  left: 8px;
  animation-delay: 0s;
}

.lds-grid div:nth-child(2) {
  top: 8px;
  left: 32px;
  animation-delay: -0.4s;
}

.lds-grid div:nth-child(3) {
  top: 8px;
  left: 56px;
  animation-delay: -0.8s;
}

.lds-grid div:nth-child(4) {
  top: 32px;
  left: 8px;
  animation-delay: -0.4s;
}

.lds-grid div:nth-child(5) {
  top: 32px;
  left: 32px;
  animation-delay: -0.8s;
}

.lds-grid div:nth-child(6) {
  top: 32px;
  left: 56px;
  animation-delay: -1.2s;
}

.lds-grid div:nth-child(7) {
  top: 56px;
  left: 8px;
  animation-delay: -0.8s;
}

.lds-grid div:nth-child(8) {
  top: 56px;
  left: 32px;
  animation-delay: -1.2s;
}

.lds-grid div:nth-child(9) {
  top: 56px;
  left: 56px;
  animation-delay: -1.6s;
}

@keyframes lds-grid {
  0%, 100% {
    opacity: 0.5;
    transform: scale(1);
  }
  50% {
    opacity: 0.25;
    transform: scale(0.75);
  }

}


.lds-spinner {
  color: rgba(255, 255, 255, 0.2);
  display: inline-block;
  position: relative;
  width: 80px;
  height: 80px;
  zoom: 0.2;
}

.lds-spinner div {
  transform-origin: 40px 40px;
  animation: lds-spinner 1.2s linear infinite;
}

.lds-spinner div:after {
  content: " ";
  display: block;
  position: absolute;
  top: 3px;
  left: 37px;
  width: 6px;
  height: 18px;
  border-radius: 20%;
  background: rgba(255, 255, 255, 0.4);
}

.lds-spinner div:nth-child(1) {
  transform: rotate(0deg);
  animation-delay: -1.1s;
}

.lds-spinner div:nth-child(2) {
  transform: rotate(30deg);
  animation-delay: -1s;
}

.lds-spinner div:nth-child(3) {
  transform: rotate(60deg);
  animation-delay: -0.9s;
}

.lds-spinner div:nth-child(4) {
  transform: rotate(90deg);
  animation-delay: -0.8s;
}

.lds-spinner div:nth-child(5) {
  transform: rotate(120deg);
  animation-delay: -0.7s;
}

.lds-spinner div:nth-child(6) {
  transform: rotate(150deg);
  animation-delay: -0.6s;
}

.lds-spinner div:nth-child(7) {
  transform: rotate(180deg);
  animation-delay: -0.5s;
}

.lds-spinner div:nth-child(8) {
  transform: rotate(210deg);
  animation-delay: -0.4s;
}

.lds-spinner div:nth-child(9) {
  transform: rotate(240deg);
  animation-delay: -0.3s;
}

.lds-spinner div:nth-child(10) {
  transform: rotate(270deg);
  animation-delay: -0.2s;
}

.lds-spinner div:nth-child(11) {
  transform: rotate(300deg);
  animation-delay: -0.1s;
}

.lds-spinner div:nth-child(12) {
  transform: rotate(330deg);
  animation-delay: 0s;
}

@keyframes lds-spinner {
  0% {
    opacity: 1;
  }
  100% {
    opacity: 0;
  }
}
</style>
