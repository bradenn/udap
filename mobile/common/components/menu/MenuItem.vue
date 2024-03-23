<!-- Copyright (c) 2023 Braden Nicholson -->

<script setup></script>

<template>
  <div class="menu-slot subplot">
    <slot></slot>
  </div>

</template>

<style lang="scss" scoped>

$border-color: rgba(255, 255, 255, 0.0625);
$border-radius: 0.5rem;

.menu-slot:not(:last-of-type):before {
  //left: 1rem;
  content: ' ';
  position: absolute;
  margin-top: calc(4rem + 1.66px);
  margin-left: -4.7rem;
  width: 75%;
  height: 2px;
  background-color: $border-color;
  z-index: -4 !important;
}

.menu-slot {
  cursor: none !important;
  user-select: none;
  height: 4rem;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 0 !important;
  border: 1px solid $border-color;
  margin: 2px;
}

.menu-slot:first-of-type {
  border-top: 1px solid $border-color !important;
  border-radius: $border-radius $border-radius 0 0 !important;
}

.menu-slot:last-of-type {
  border-radius: 0 0 $border-radius $border-radius !important;
  border-bottom: 1px solid $border-color !important;
}

.menu-slot:not(:last-of-type:first-of-type) {
  border-top: none;
  border-bottom: none;
}

</style>