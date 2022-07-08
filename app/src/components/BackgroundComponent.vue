<script setup lang="ts">
import logo from "@/assets/logo.svg";
</script>

<template>
  <div class="background-container">
    <div class="background-transformed"></div>
  </div>
  <div class="banner">
    <img :src="logo" width="60" />
    <span class="banner-text text-4xl font-bold mx-5">Kecp</span>
  </div>
</template>

<style scoped>
.background-container {
  position: absolute;
  top: 0;
  left: 0;
  min-height: 100vh;
  width: 100vw;
  background: #f6f9fc;
  z-index: -100;
}

.banner {
  position: absolute;
  top: 22vh;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: -80;
}

@media (max-height: 768px) {
  .banner {
    top: 16vh;
  }
}

.banner-text {
  color: #ffffff;
}

.background-transformed {
  background: #525f7f;
  transform: skewY(-10deg);
  height: 70vh;
  width: 100vw;
  transform-origin: 0% 100%;
  z-index: -99;
}
</style>
